package tvd

// Denoise1D applies Total variation denoising to Y and stores it to X
// for the given lambda.
func Denoise1D(Y []float64, lambda float64) []float64 {
	N := len(Y)
	X := make([]float64, N)
	var k, k0, kz, kf int
	vmin := Y[0] - lambda
	vmax := Y[0] + lambda
	umin := lambda
	umax := -lambda

	for k < N {
		if k == N-1 {
			X[k] = vmin + umin
			break
		}
		y := Y[k]
		ynext := Y[k+1]
		switch {
		case ynext < vmin-lambda-umin:
			for i := k0; i < kf+1; i++ {
				X[i] = vmin
			}
			k, k0, kz, kf = kf+1, kf+1, kf+1, kf+1
			vmin = y
			vmax = y + 2*lambda
			umin = lambda
			umax = -lambda

		case ynext > vmax+lambda-umax:
			for i := k0; i < kz+1; i++ {
				X[i] = vmax
			}
			k, k0, kz, kf = kz+1, kz+1, kz+1, kz+1
			vmin = y - 2*lambda
			vmax = y
			umin = lambda
			umax = -lambda

		default:
			k++
			y = ynext
			umin += y - vmin
			umax += y - vmax
			if umin >= lambda {
				vmin += (umin - lambda) / float64(k-k0+1)
				umin = lambda
				kf = k
			}
			if umax <= -lambda {
				vmax += (umax + lambda) / float64(k-k0+1)
				umax = -lambda
				kz = k
			}
		}
		if k != N-1 {
			continue
		}

		if umin < 0 {
			for i := k0; i < kf+1; i++ {
				X[i] = vmin
			}
			k, k0, kf = kf+1, kf+1, kf+1
			y = Y[k]
			vmin = y
			umin = lambda
			umax = y + lambda - vmax

		} else if umax > 0 {
			for i := k0; i < kz+1; i++ {
				X[i] = vmax
			}
			k, k0, kz = kz+1, kz+1, kz+1
			y = Y[k]
			vmax = y
			umax = -lambda
			umin = y - lambda - vmin
		} else {
			for i := k0; i < k+1; i++ {
				X[i] = vmin + umin/float64(k-k0+1)
			}
			break
		}
	}
	return X
}
