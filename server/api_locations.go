package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/etzba/pggo/dat"
	"github.com/etzba/pggo/wire"
)

func (s *Server) getLocations() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		s.Logger.Info("Server go request" + " method: " + r.Method + " uri: " + r.RequestURI)
		now := time.Now()
		defer s.shipper.Collect(now, r)
		locations, err := s.Database.GetAllLocationDetails()
		if err != nil {
			s.Logger.Error("Failed to get all locations", err)
			s.Respoder.SendError(w, err)
			return
		}

		for _, l := range locations {
			_, err := fmt.Fprintf(w, "Location: %+v\n", l)
			if err != nil {
				s.Logger.Error(fmt.Sprintf("Failed to print location %s", l.Name), err)
			}
		}
		s.Respoder.SendOK(w, locations)
	}
}

func (s *Server) getLocationById() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		s.Logger.Info("Server go request" + " method: " + r.Method + " uri: " + r.RequestURI)
		now := time.Now()
		defer s.shipper.Collect(now, r)
		idStr, _ := strings.CutPrefix(r.URL.Path, "/locations/")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			s.Logger.Error("Failed to convert string to integer", err)
			s.Respoder.SendError(w, err)
			return
		}

		location, err := s.Database.GetLocationDetailsByID(id)
		if err != nil {
			s.Logger.Error("Failed to get location by id", err)
			s.Respoder.SendError(w, err)
			return
		}

		_, err = fmt.Fprintf(w, "Location: %+v\n", location)
		if err != nil {
			s.Logger.Error(fmt.Sprintf("Failed to print location %s", location.Name), err)
		}
		s.Respoder.SendOK(w, location)
	}
}

func (s *Server) addLocation() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		s.Logger.Info("Server go request" + " method: " + r.Method + " uri: " + r.RequestURI)
		now := time.Now()
		defer s.shipper.Collect(now, r)
		loc := wire.Location{}
		if err := json.NewDecoder(r.Body).Decode(&loc); err != nil {
			s.Logger.Error("Failed to create new decoder", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		location := dat.Location{
			Name:       loc.Name,
			Address:    loc.Address,
			Longtitude: loc.Longtitude,
			Latitude:   loc.Latitude,
		}

		if err := s.Database.InsertLocationIntoDatabase(location); err != nil {
			s.Logger.Error("Failed to insert location", err)
			s.Respoder.SendError(w, err)
			return
		}

		_, err := fmt.Fprintf(w, "Location added: %+v\n", location)
		if err != nil {
			s.Logger.Error(fmt.Sprintf("Failed to print location %s", location.Name), err)
		}
		s.Respoder.SendOK(w, location)
	}
}

func (s *Server) updateLocation() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		s.Logger.Info("Server go request" + " method: " + r.Method + " uri: " + r.RequestURI)
		now := time.Now()
		defer s.shipper.Collect(now, r)
		idStr, _ := strings.CutPrefix(r.URL.Path, "/locations/")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			s.Logger.Error("Failed to convert string to integer", err)
			s.Respoder.SendError(w, err)
			return
		}

		loc := wire.Location{}
		if err := json.NewDecoder(r.Body).Decode(&loc); err != nil {
			s.Logger.Error("Failed to create new decoder", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		location := dat.Location{
			Name:       loc.Name,
			Address:    loc.Address,
			Longtitude: loc.Longtitude,
			Latitude:   loc.Latitude,
		}

		if err := s.Database.UpdateLocationDetails(id, location); err != nil {
			s.Logger.Error("Failed to update location", err)
			s.Respoder.SendError(w, err)
			return
		}

		_, err = fmt.Fprintf(w, "Location updated: %+v\n", location)
		if err != nil {
			s.Logger.Error(fmt.Sprintf("Failed to print location %s", location.Name), err)
		}
		s.Respoder.SendOK(w, location)
	}
}

func (s *Server) deleteLocationById() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		s.Logger.Info("Server go request" + " method: " + r.Method + " uri: " + r.RequestURI)
		now := time.Now()
		defer s.shipper.Collect(now, r)
		idStr, _ := strings.CutPrefix(r.URL.Path, "/locations/")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			s.Logger.Error("Failed to convert string to integer", err)
			s.Respoder.SendError(w, err)
			return
		}

		if err := s.Database.DeleteLocationFromDatabase(id); err != nil {
			s.Logger.Error("Failed to delete location by id", err)
			s.Respoder.SendError(w, err)
			return
		}

		_, err = fmt.Fprintf(w, "Location id deleted: %+v\n", id)
		if err != nil {
			s.Logger.Error(fmt.Sprintf("Failed to print location %d", id), err)
		}
		s.Respoder.SendOK(w, id)
	}
}
