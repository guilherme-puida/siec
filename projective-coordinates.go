package siec

import (
	"math/big"
)

func affineToProjective(x, y *big.Int) (X, Y, Z *big.Int) {
	return new(big.Int).Set(x), new(big.Int).Set(y), big.NewInt(1)
}

// copies from https://golang.org/src/crypto/elliptic/elliptic.go
func projectiveToAffine(x, y, z *big.Int) (X, Y *big.Int) {
	curve := SIEC255()
	if z.Sign() == 0 {
		return new(big.Int), new(big.Int)
	}
	zinv := new(big.Int).ModInverse(z, curve.P)
	zinvsq := new(big.Int).Mul(zinv, zinv)
	X = new(big.Int).Mul(x, zinvsq)
	X.Mod(X, curve.P)
	Y = new(big.Int).Mul(y, zinvsq.Mul(zinvsq, zinv))
	Y.Mod(Y, curve.P)
	return
}

// http://hyperelliptic.org/EFD/g1p/auto-shortw-jacobian-0.html#addition-add-2001-b
func add2007bl(X1, Y1, Z1, X2, Y2, Z2 *big.Int) (X3, Y3, Z3 *big.Int) {
	w := new(big.Int)
	ww := new(big.Int)
	curve := SIEC255()
	// Z1Z1 = Z1^2
	Z1Z1 := new(big.Int).Mul(Z1, Z1)
	// Z2Z2 = Z2^2
	Z2Z2 := new(big.Int).Mul(Z2, Z2)
	// U1 = X1*Z2Z2
	U1 := new(big.Int).Mul(X1, Z2Z2)
	// U2 = X2*Z1Z1
	U2 := new(big.Int).Mul(X2, Z1Z1)
	// S1 = Y1*Z2*Z2Z2
	S1 := new(big.Int).Mul(Y1, w.Mul(Z2, Z2Z2))
	// S2 = Y2*Z1*Z1Z1
	S2 := new(big.Int).Mul(Y2, w.Mul(Z1, Z1Z1))
	// H = U2-U1
	H := new(big.Int).Sub(U2, U1)
	// I = (2*H)^2
	I := new(big.Int).Exp(w.Lsh(H, 1), two, curve.P)
	// J = H*I
	J := new(big.Int).Mul(H, I)
	// r = 2*(S2-S1)
	r := new(big.Int).Lsh(w.Sub(S2, S1), 1)
	// V = U1*I
	V := new(big.Int).Mul(U1, I)
	// X3 = r^2-J-2*V
	X3 = new(big.Int).Sub(
		w.Mul(r, r),
		ww.Add(
			J,
			ww.Lsh(V, 1),
		),
	)
	// Y3 = r*(V-X3)-2*S1*J
	Y3 = new(big.Int).Sub(
		w.Mul(r, w.Sub(V, X3)),
		ww.Lsh(ww.Mul(S1, J), 1),
	)
	// Z3 = ((Z1+Z2)^2-Z1Z1-Z2Z2)*H
	Z3 = new(big.Int).Mul(
		w.Sub(
			w.Mul(w.Add(Z1, Z2), w),
			ww.Add(Z1Z1, Z2Z2),
		),
		H,
	)
	return
}

func dbl2009l(X1, Y1, Z1 *big.Int) (X3, Y3, Z3 *big.Int) {
	w := new(big.Int)
	m := new(big.Int)
	curve := SIEC255()
	// A = X1^2
	A := new(big.Int).Mul(X1, X1)
	A.Mod(A, curve.P)
	// B = Y1^2
	B := new(big.Int).Mul(Y1, Y1)
	B.Mod(B, curve.P)
	// C = B^2
	C := new(big.Int).Mul(B, B)
	C.Mod(C, curve.P)
	// D = 2*((X1+B)^2-A-C)
	w.Add(X1, B)
	D := new(big.Int).Lsh(w.Sub(w.Mul(w, w), m.Add(A, C)), 1)
	D.Mod(D, curve.P)
	// E = 3*A
	E := new(big.Int).Mul(three, A)
	E.Mod(E, curve.P)
	// F = E^2
	F := new(big.Int).Mul(E, E)
	F.Mod(F, curve.P)
	// X3 = F-2*D
	X3 = new(big.Int).Sub(F, w.Lsh(D, 1))
	X3.Mod(X3, curve.P)
	// Y3 = E*(D-X3)-8*C
	Y3 = new(big.Int).Sub(w.Mul(E, w.Sub(D, X3)), m.Mul(eight, C))
	Y3.Mod(Y3, curve.P)
	// Z3 = 2*Y1*Z1
	Z3 = w.Lsh(w.Mul(Z1, Y1), 1)
	Z3.Mod(Z3, curve.P)
	return
}
