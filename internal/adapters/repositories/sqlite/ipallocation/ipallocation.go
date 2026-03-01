package ipallocation

import (
	"context"
	"database/sql"
	"fmt"
	"net/netip"

	sqlc "github.com/thebeyond-net/node-agent/internal/adapters/repositories/sqlite/generated"
)

var ErrIPPoolExhausted = fmt.Errorf("ip pool exhausted")

type Repository struct {
	db      *sql.DB
	queries *sqlc.Queries
}

func NewRepository(db *sql.DB) (*Repository, error) {
	queries := sqlc.New(db)
	return &Repository{db, queries}, nil
}

func (r *Repository) Allocate(ctx context.Context, prefix netip.Prefix) (netip.Addr, error) {
	if ip, err := r.tryReclaim(ctx, prefix); err == nil {
		_ = r.updateNextCandidate(ctx, prefix, ip.Next())
		return ip, nil
	}

	candidate := r.getNextCandidate(ctx, prefix)
	start := candidate

	for prefix.Contains(candidate) {
		allocated, err := r.isAllocated(ctx, candidate)
		if err != nil {
			return netip.Addr{}, fmt.Errorf("check allocation: %w", err)
		}
		if !allocated {
			_ = r.updateNextCandidate(ctx, prefix, candidate.Next())
			return candidate, nil
		}
		candidate = candidate.Next()
	}

	candidate = prefix.Addr().Next().Next()
	for candidate.Compare(start) < 0 && prefix.Contains(candidate) {
		allocated, err := r.isAllocated(ctx, candidate)
		if err != nil {
			return netip.Addr{}, fmt.Errorf("check allocation: %w", err)
		}
		if !allocated {
			_ = r.updateNextCandidate(ctx, prefix, candidate.Next())
			return candidate, nil
		}
		candidate = candidate.Next()
	}

	return netip.Addr{}, fmt.Errorf("%w in %s", ErrIPPoolExhausted, prefix)
}

func (r *Repository) tryReclaim(ctx context.Context, prefix netip.Prefix) (netip.Addr, error) {
	ip, err := r.queries.GetReleasedIP(ctx, prefix.String())
	if err != nil {
		return netip.Addr{}, err
	}
	return netip.MustParseAddr(ip), r.queries.ReleaseIP(ctx, ip)
}

func (r *Repository) getNextCandidate(ctx context.Context, prefix netip.Prefix) netip.Addr {
	candidate, err := r.queries.GetNextCandidate(ctx, prefix.String())
	if err != nil {
		return prefix.Addr().Next().Next()
	}
	return netip.MustParseAddr(candidate)
}

func (r *Repository) updateNextCandidate(ctx context.Context, prefix netip.Prefix, next netip.Addr) error {
	if err := r.queries.UpdateNextCandidate(ctx, sqlc.UpdateNextCandidateParams{
		Prefix:        prefix.String(),
		NextCandidate: next.String(),
	}); err != nil {
		return fmt.Errorf("update next candidate: %w", err)
	}
	return nil
}

func (r *Repository) isAllocated(ctx context.Context, ip netip.Addr) (bool, error) {
	cnt, err := r.queries.IsAllocated(ctx, ip.String())
	if err != nil {
		return false, err
	}
	return cnt > 0, nil
}

func (r *Repository) Reserve(ctx context.Context, ip netip.Addr, publicKey string, prefix netip.Prefix) error {
	return r.queries.Reserve(ctx, sqlc.ReserveParams{
		Ip:     ip.String(),
		Pubkey: publicKey,
		Prefix: prefix.String(),
	})
}

func (r *Repository) Release(ctx context.Context, ip netip.Addr) error {
	return r.queries.Release(ctx, ip.String())
}

func (r *Repository) ReleaseByPublicKey(ctx context.Context, publicKey string) error {
	return r.queries.ReleaseByPublicKey(ctx, publicKey)
}
