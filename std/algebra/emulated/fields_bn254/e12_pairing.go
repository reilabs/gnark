package fields_bn254

import (
	"github.com/consensys/gnark/std/math/emulated"
)

func (e Ext12) nSquareTorus(z *E6, n int) *E6 {
	for i := 0; i < n; i++ {
		z = e.SquareTorus(z)
	}
	return z
}

// Exponentiation by the seed t=4965661367192848881
// The computations are performed on E6 compressed form using Torus-based arithmetic.
func (e Ext12) ExptTorus(x *E6) *E6 {
	// Expt computation is derived from the addition chain:
	//
	//	_10     = 2*1
	//	_100    = 2*_10
	//	_1000   = 2*_100
	//	_10000  = 2*_1000
	//	_10001  = 1 + _10000
	//	_10011  = _10 + _10001
	//	_10100  = 1 + _10011
	//	_11001  = _1000 + _10001
	//	_100010 = 2*_10001
	//	_100111 = _10011 + _10100
	//	_101001 = _10 + _100111
	//	i27     = (_100010 << 6 + _100 + _11001) << 7 + _11001
	//	i44     = (i27 << 8 + _101001 + _10) << 6 + _10001
	//	i70     = ((i44 << 8 + _101001) << 6 + _101001) << 10
	//	return    (_100111 + i70) << 6 + _101001 + _1000
	//
	// Operations: 62 squares 17 multiplies
	//
	// Generated by github.com/mmcloughlin/addchain v0.4.0.

	t3 := e.SquareTorus(x)
	t5 := e.SquareTorus(t3)
	result := e.SquareTorus(t5)
	t0 := e.SquareTorus(result)
	t2 := e.MulTorus(x, t0)
	t0 = e.MulTorus(t3, t2)
	t1 := e.MulTorus(x, t0)
	t4 := e.MulTorus(result, t2)
	t6 := e.SquareTorus(t2)
	t1 = e.MulTorus(t0, t1)
	t0 = e.MulTorus(t3, t1)
	t6 = e.nSquareTorus(t6, 6)
	t5 = e.MulTorus(t5, t6)
	t5 = e.MulTorus(t4, t5)
	t5 = e.nSquareTorus(t5, 7)
	t4 = e.MulTorus(t4, t5)
	t4 = e.nSquareTorus(t4, 8)
	t4 = e.MulTorus(t0, t4)
	t3 = e.MulTorus(t3, t4)
	t3 = e.nSquareTorus(t3, 6)
	t2 = e.MulTorus(t2, t3)
	t2 = e.nSquareTorus(t2, 8)
	t2 = e.MulTorus(t0, t2)
	t2 = e.nSquareTorus(t2, 6)
	t2 = e.MulTorus(t0, t2)
	t2 = e.nSquareTorus(t2, 10)
	t1 = e.MulTorus(t1, t2)
	t1 = e.nSquareTorus(t1, 6)
	t0 = e.MulTorus(t0, t1)
	z := e.MulTorus(result, t0)
	return z
}

// Square034 squares an E12 sparse element of the form
//
//	E12{
//		C0: E6{B0: 1, B1: 0, B2: 0},
//		C1: E6{B0: c3, B1: c4, B2: 0},
//	}
func (e *Ext12) Square034(x *E12) *E12 {
	c0 := E6{
		B0: *e.Ext2.Sub(&x.C0.B0, &x.C1.B0),
		B1: *e.Ext2.Neg(&x.C1.B1),
		B2: *e.Ext2.Zero(),
	}

	c3 := &E6{
		B0: x.C0.B0,
		B1: *e.Ext2.Neg(&x.C1.B0),
		B2: *e.Ext2.Neg(&x.C1.B1),
	}

	c2 := E6{
		B0: x.C1.B0,
		B1: x.C1.B1,
		B2: *e.Ext2.Zero(),
	}
	c3 = e.MulBy01(c3, &c0.B0, &c0.B1)
	c3 = e.Ext6.Add(c3, &c2)

	var z E12
	z.C1.B0 = *e.Ext2.Add(&c2.B0, &c2.B0)
	z.C1.B1 = *e.Ext2.Add(&c2.B1, &c2.B1)

	z.C0.B0 = c3.B0
	z.C0.B1 = *e.Ext2.Add(&c3.B1, &c2.B0)
	z.C0.B2 = *e.Ext2.Add(&c3.B2, &c2.B1)

	return &z
}

// MulBy034 multiplies z by an E12 sparse element of the form
//
//	E12{
//		C0: E6{B0: 1, B1: 0, B2: 0},
//		C1: E6{B0: c3, B1: c4, B2: 0},
//	}
func (e *Ext12) MulBy034(z *E12, c3, c4 *E2) *E12 {

	a := z.C0
	b := e.MulBy01(&z.C1, c3, c4)
	c3 = e.Ext2.Add(e.Ext2.One(), c3)
	d := e.Ext6.Add(&z.C0, &z.C1)
	d = e.MulBy01(d, c3, c4)

	zC1 := e.Ext6.Add(&a, b)
	zC1 = e.Ext6.Neg(zC1)
	zC1 = e.Ext6.Add(zC1, d)
	zC0 := e.Ext6.MulByNonResidue(b)
	zC0 = e.Ext6.Add(zC0, &a)

	return &E12{
		C0: *zC0,
		C1: *zC1,
	}
}

//	multiplies two E12 sparse element of the form:
//
//	E12{
//		C0: E6{B0: 1, B1: 0, B2: 0},
//		C1: E6{B0: c3, B1: c4, B2: 0},
//	}
//
// and
//
//	E12{
//		C0: E6{B0: 1, B1: 0, B2: 0},
//		C1: E6{B0: d3, B1: d4, B2: 0},
//	}
func (e *Ext12) Mul034By034(d3, d4, c3, c4 *E2) [5]*E2 {
	x3 := e.Ext2.Mul(c3, d3)
	x4 := e.Ext2.Mul(c4, d4)
	x04 := e.Ext2.Add(c4, d4)
	x03 := e.Ext2.Add(c3, d3)
	tmp := e.Ext2.Add(c3, c4)
	x34 := e.Ext2.Add(d3, d4)
	x34 = e.Ext2.Mul(x34, tmp)
	x34 = e.Ext2.Sub(x34, x3)
	x34 = e.Ext2.Sub(x34, x4)

	zC0B0 := e.Ext2.MulByNonResidue(x4)
	zC0B0 = e.Ext2.Add(zC0B0, e.Ext2.One())
	zC0B1 := x3
	zC0B2 := x34
	zC1B0 := x03
	zC1B1 := x04

	return [5]*E2{zC0B0, zC0B1, zC0B2, zC1B0, zC1B1}
}

// MulBy01234 multiplies z by an E12 sparse element of the form
//
//	E12{
//		C0: E6{B0: c0, B1: c1, B2: c2},
//		C1: E6{B0: c3, B1: c4, B2: 0},
//	}
func (e *Ext12) MulBy01234(z *E12, x [5]*E2) *E12 {
	c0 := &E6{B0: *x[0], B1: *x[1], B2: *x[2]}
	c1 := &E6{B0: *x[3], B1: *x[4], B2: *e.Ext2.Zero()}
	a := e.Ext6.Add(&z.C0, &z.C1)
	b := e.Ext6.Add(c0, c1)
	a = e.Ext6.Mul(a, b)
	b = e.Ext6.Mul(&z.C0, c0)
	c := e.Ext6.MulBy01(&z.C1, x[3], x[4])
	z1 := e.Ext6.Sub(a, b)
	z1 = e.Ext6.Sub(z1, c)
	z0 := e.Ext6.MulByNonResidue(c)
	z0 = e.Ext6.Add(z0, b)
	return &E12{
		C0: *z0,
		C1: *z1,
	}
}

//	multiplies two E12 sparse element of the form:
//
//	E12{
//		C0: E6{B0: x0, B1: x1, B2: x2},
//		C1: E6{B0: x3, B1: x4, B2: 0},
//	}
//
// and
//
//	E12{
//		C0: E6{B0: 1, B1: 0, B2: 0},
//		C1: E6{B0: z3, B1: z4, B2: 0},
//	}
func (e *Ext12) Mul01234By034(x [5]*E2, z3, z4 *E2) *E12 {
	c0 := &E6{B0: *x[0], B1: *x[1], B2: *x[2]}
	c1 := &E6{B0: *x[3], B1: *x[4], B2: *e.Ext2.Zero()}
	a := e.Ext6.Add(e.Ext6.One(), &E6{B0: *z3, B1: *z4, B2: *e.Ext2.Zero()})
	b := e.Ext6.Add(c0, c1)
	a = e.Ext6.Mul(a, b)
	c := e.Ext6.Mul01By01(z3, z4, x[3], x[4])
	z1 := e.Ext6.Sub(a, c0)
	z1 = e.Ext6.Sub(z1, c)
	z0 := e.Ext6.MulByNonResidue(c)
	z0 = e.Ext6.Add(z0, c0)
	return &E12{
		C0: *z0,
		C1: *z1,
	}
}

// Torus-based arithmetic:
//
// After the easy part of the final exponentiation the elements are in a proper
// subgroup of Fpk (E12) that coincides with some algebraic tori. The elements
// are in the torus Tk(Fp) and thus in each torus Tk/d(Fp^d) for d|k, d≠k.  We
// take d=6. So the elements are in T2(Fp6).
// Let G_{q,2} = {m ∈ Fq^2 | m^(q+1) = 1} where q = p^6.
// When m.C1 = 0, then m.C0 must be 1 or −1.
//
// We recall the tower construction:
//
//	𝔽p²[u] = 𝔽p/u²+1
//	𝔽p⁶[v] = 𝔽p²/v³-9-u
//	𝔽p¹²[w] = 𝔽p⁶/w²-v

// CompressTorus compresses x ∈ E12 to (x.C0 + 1)/x.C1 ∈ E6
func (e Ext12) CompressTorus(x *E12) *E6 {
	// x ∈ G_{q,2} \ {-1,1}
	y := e.Ext6.Add(&x.C0, e.Ext6.One())
	y = e.Ext6.DivUnchecked(y, &x.C1)
	return y
}

// DecompressTorus decompresses y ∈ E6 to (y+w)/(y-w) ∈ E12
func (e Ext12) DecompressTorus(y *E6) *E12 {
	var n, d E12
	one := e.Ext6.One()
	n.C0 = *y
	n.C1 = *one
	d.C0 = *y
	d.C1 = *e.Ext6.Neg(one)

	x := e.DivUnchecked(&n, &d)
	return x
}

// MulTorus multiplies two compressed elements y1, y2 ∈ E6
// and returns (y1 * y2 + v)/(y1 + y2)
// N.B.: we use MulTorus in the final exponentiation throughout y1 ≠ -y2 always.
func (e Ext12) MulTorus(y1, y2 *E6) *E6 {
	n := e.Ext6.Mul(y1, y2)
	n = &E6{
		B0: n.B0,
		B1: *e.Ext2.Add(&n.B1, e.Ext2.One()),
		B2: n.B2,
	}
	d := e.Ext6.Add(y1, y2)
	y3 := e.Ext6.DivUnchecked(n, d)
	return y3
}

// InverseTorus inverses a compressed elements y ∈ E6
// and returns -y
func (e Ext12) InverseTorus(y *E6) *E6 {
	return e.Ext6.Neg(y)
}

// SquareTorus squares a compressed elements y ∈ E6
// and returns (y + v/y)/2
//
// It uses a hint to verify that (2x-y)y = v saving one E6 AssertIsEqual.
func (e Ext12) SquareTorus(y *E6) *E6 {
	res, err := e.fp.NewHint(squareTorusHint, 6, &y.B0.A0, &y.B0.A1, &y.B1.A0, &y.B1.A1, &y.B2.A0, &y.B2.A1)
	if err != nil {
		// err is non-nil only for invalid number of inputs
		panic(err)
	}

	sq := E6{
		B0: E2{A0: *res[0], A1: *res[1]},
		B1: E2{A0: *res[2], A1: *res[3]},
		B2: E2{A0: *res[4], A1: *res[5]},
	}

	// v = (2x-y)y
	v := e.Ext6.Double(&sq)
	v = e.Ext6.Sub(v, y)
	v = e.Ext6.Mul(v, y)

	_v := E6{B0: *e.Ext2.Zero(), B1: *e.Ext2.One(), B2: *e.Ext2.Zero()}
	e.Ext6.AssertIsEqual(v, &_v)

	return &sq

}

// FrobeniusTorus raises a compressed elements y ∈ E6 to the modulus p
// and returns y^p / v^((p-1)/2)
func (e Ext12) FrobeniusTorus(y *E6) *E6 {
	t0 := e.Ext2.Conjugate(&y.B0)
	t1 := e.Ext2.Conjugate(&y.B1)
	t2 := e.Ext2.Conjugate(&y.B2)
	t1 = e.Ext2.MulByNonResidue1Power2(t1)
	t2 = e.Ext2.MulByNonResidue1Power4(t2)

	v0 := E2{emulated.ValueOf[emulated.BN254Fp]("18566938241244942414004596690298913868373833782006617400804628704885040364344"), emulated.ValueOf[emulated.BN254Fp]("5722266937896532885780051958958348231143373700109372999374820235121374419868")}
	res := &E6{B0: *t0, B1: *t1, B2: *t2}
	res = e.Ext6.MulBy0(res, &v0)

	return res
}

// FrobeniusSquareTorus raises a compressed elements y ∈ E6 to the square modulus p^2
// and returns y^(p^2) / v^((p^2-1)/2)
func (e Ext12) FrobeniusSquareTorus(y *E6) *E6 {
	v0 := emulated.ValueOf[emulated.BN254Fp]("2203960485148121921418603742825762020974279258880205651967")
	t0 := e.Ext2.MulByElement(&y.B0, &v0)
	t1 := e.Ext2.MulByNonResidue2Power2(&y.B1)
	t1 = e.Ext2.MulByElement(t1, &v0)
	t2 := e.Ext2.MulByNonResidue2Power4(&y.B2)
	t2 = e.Ext2.MulByElement(t2, &v0)

	return &E6{B0: *t0, B1: *t1, B2: *t2}
}

// FrobeniusCubeTorus raises a compressed elements y ∈ E6 to the cube modulus p^3
// and returns y^(p^3) / v^((p^3-1)/2)
func (e Ext12) FrobeniusCubeTorus(y *E6) *E6 {
	t0 := e.Ext2.Conjugate(&y.B0)
	t1 := e.Ext2.Conjugate(&y.B1)
	t2 := e.Ext2.Conjugate(&y.B2)
	t1 = e.Ext2.MulByNonResidue3Power2(t1)
	t2 = e.Ext2.MulByNonResidue3Power4(t2)

	v0 := E2{emulated.ValueOf[emulated.BN254Fp]("10190819375481120917420622822672549775783927716138318623895010788866272024264"), emulated.ValueOf[emulated.BN254Fp]("303847389135065887422783454877609941456349188919719272345083954437860409601")}
	res := &E6{B0: *t0, B1: *t1, B2: *t2}
	res = e.Ext6.MulBy0(res, &v0)

	return res
}