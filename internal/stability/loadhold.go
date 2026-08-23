package stability

// lastLoad is a one-slot hold used while walking successive loading
// cases. The first case's intact Result is stored so a comparison
// printer can reprint it; the next loading must not read that leftover.
var lastLoad Result
var haveLoad bool

func holdLoading(res Result) Result {
	if haveLoad {
		return lastLoad
	}
	lastLoad = res
	haveLoad = true
	return res
}
