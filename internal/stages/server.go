package stages

type Server interface {
	GetOccupyingCount() int
	AddCar(car *Car) bool
	Serve() <-chan *Car
	CloseInbound()
}

func FindMin(servers []*Server) *Server {
	first := servers[0]

	for _, server := range servers[1:] {
		if (*server).GetOccupyingCount() < (*first).GetOccupyingCount() {
			first = server
		}
	}

	return first
}
