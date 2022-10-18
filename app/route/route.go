package route

import "errors"

//? struct: generates a data structure
//? So we can works using this structure
type Route struct {
	ID string `json: "routeId"`
	ClientId string `json: "clientID"`
	Positions []Position `json: "positions"`
}

type Position struct {
	Lat float64 `json: "lat"` 		//? Latitude
	Long float64 `json: "long"`		//? Longitude
}

//? What our backend needs
type PartialRoutePosition struct {
	ID string `json: "routeId"` //? Passing this string will create the key when converted to json
	ClientID string `json: "clientID"`
	Position []float64 `json: "positions"`
	Finished bool `json: finished`
}

func(r *Route) LoadPositions() error {
	if r.ID == "" {
		return errors.New("route id not informed")
	}

	f, err := os.Open("destinations/" + r.ID + ".txt")
	if err != nil {
		return err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		data := string.Split(scanner.Text(), ",")
		lat, err := strconv.ParseFloat(data[0], 64)
		if err != nil {
			return nil
		}
		long, err := strconv.ParseFloat(data[1], 64)
		if err != nil {
			return nil
		}
		r.Positions = append(r.Positions, Position{
			Lat: lat,
			Long: long
		})
	}
	return nill
}


//? We need to generate a JSON with a
//? list of positions to send to Kafka
func (r *Route) ExportJsonPOsitions() ([]string, error) {
	var route PartialRoutePosition
	var result []string
	total := len(r.Positions)

	for k, v := range r.Positions {
		route.ID = r.ID
		route.ClientId = r.ClientID
		route.Position = []float64{v.Lat, v.Long}
		route.Finished = false

		//? When finished the Positions we need to
		//? send for the backend telling that our
		//? delivery track its over
		if total-1 == k {
			route.Finished = true
		}

		jsonRoute, err := json.Marshal(route)
		if err != nil {
			return nil, err
		}

		result = append(result, string(jsonRoute))
	}

	return result, nil
}

