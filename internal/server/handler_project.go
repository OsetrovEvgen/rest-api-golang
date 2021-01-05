package server

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/osetr/rest-api/internal/model"
)

// Add new project
func (s *APIServer) createProject() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		tx, err := s.Store.DB.Begin()
		req := &model.Project{}
		if err != nil {
			s.errorresp(w, r, http.StatusInternalServerError, err)
			return
		}

		if err = json.NewDecoder(r.Body).Decode(req); err != nil {
			s.errorresp(w, r, http.StatusBadRequest, err)
			return
		}
		if err = req.Validate(); err != nil {
			s.errorresp(w, r, http.StatusBadRequest, err)
			return
		}

		p := &model.Project{
			ID:          genID(projectID),
			Name:        req.Name,
			Description: req.Description,
		}
		position := 0
		name := "Default column"
		c := &model.Column{
			ID:        genID(columnID),
			ProjectID: p.ID,
			Name:      &name,
			Position:  &position,
		}

		if _, err = tx.Exec(
			"INSERT INTO projects (id, name, description) VALUES ($1, $2, $3)",
			p.ID,
			p.Name,
			p.Description,
		); err != nil {
			s.errorresp(w, r, http.StatusInternalServerError, err)
			return
		}

		if _, err = tx.Exec(
			"INSERT INTO columns (id, project_id, name, position) VALUES ($1, $2, $3, $4)",
			c.ID,
			c.ProjectID,
			c.Name,
			c.Position,
		); err != nil {
			s.errorresp(w, r, http.StatusInternalServerError, err)
			tx.Rollback()
			return
		}

		err = tx.Commit()
		if err != nil {
			s.errorresp(w, r, http.StatusInternalServerError, err)
			return
		}
		s.successresp(w, r, http.StatusCreated, p)
	}
}

// Get projects list
func (s *APIServer) getProjects() http.HandlerFunc {
	min := func(a, b int) int {
		if a < b {
			return a
		}
		return b
	}

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		query := r.URL.Query()
		pageSize := int(10)
		page := int(0)

		if _, present := query["page_size"]; present {
			ps, err := strconv.ParseInt(query.Get("page_size"), 10, 64)
			if err != nil || ps < 0 {
				s.errorresp(w, r, http.StatusBadRequest, errors.New("uncorrect page_size param"))
				return
			}
			pageSize = int(ps)
		}
		if _, present := query["page"]; present {
			p, err := strconv.ParseInt(query.Get("page"), 10, 64)
			if err != nil || p < 0 {
				s.errorresp(w, r, http.StatusBadRequest, errors.New("uncorrect paga param"))
				return
			}
			page = int(p)
		}

		rows, err := s.Store.DB.Query("SELECT * FROM projects ORDER BY name")
		if err != nil {
			s.errorresp(w, r, http.StatusInternalServerError, err)
			return
		}
		result := []model.Project{}
		for rows.Next() {
			p := &model.Project{}
			err = rows.Scan(
				&p.ID,
				&p.Name,
				&p.Description)
			result = append(result, *p)
		}

		pstart := min(page*pageSize, int(len(result)))
		pfinish := min(page*pageSize+pageSize, int(len(result)))
		s.successresp(w, r, http.StatusOK, map[string]interface{}{
			"total_count": len(result),
			"data":        result[pstart:pfinish],
		})
	}
}

// Get single project
func (s *APIServer) getProject() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		params := mux.Vars(r)

		p := &model.Project{}
		if err := s.Store.DB.QueryRow(
			"SELECT id, name, description FROM projects WHERE id = $1",
			params["id"],
		).Scan(
			&p.ID,
			&p.Name,
			&p.Description,
		); err != nil {
			s.errorresp(w, r, http.StatusNotFound, errors.New("project not found"))
			return
		}

		s.successresp(w, r, http.StatusOK, p)
	}
}

// Update project
func (s *APIServer) updateProject() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		req := &model.Project{}
		params := mux.Vars(r)

		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.errorresp(w, r, http.StatusBadRequest, err)
			return
		}
		if err := req.Validate(); err != nil {
			s.errorresp(w, r, http.StatusBadRequest, err)
			return
		}

		id := params["id"]
		p := &model.Project{
			ID:          &id,
			Name:        req.Name,
			Description: req.Description,
		}
		res, err := s.Store.DB.Exec(
			"UPDATE projects SET name = $2, description = $3 WHERE id = $1",
			p.ID,
			p.Name,
			p.Description,
		)
		if err != nil {
			s.errorresp(w, r, http.StatusInternalServerError, err)
			return
		}
		if st, _ := res.RowsAffected(); st == 0 {
			s.errorresp(w, r, http.StatusNotFound, errors.New("project not found"))
			return
		}

		s.successresp(w, r, http.StatusOK, p)
	}
}

// Patch project
func (s *APIServer) patchProject() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		req := &model.Project{}
		params := mux.Vars(r)
		var err error

		if err = json.NewDecoder(r.Body).Decode(req); err != nil {
			s.errorresp(w, r, http.StatusBadRequest, err)
			return
		}
		if err = req.ValidatePatch(); err != nil {
			s.errorresp(w, r, http.StatusBadRequest, err)
			return
		}

		id := params["id"]
		p := &model.Project{
			ID:          &id,
			Name:        req.Name,
			Description: req.Description,
		}

		var res sql.Result
		if p.Description == nil && p.Name != nil {
			res, err = s.Store.DB.Exec(
				"UPDATE projects SET name = $2 WHERE id = $1",
				p.ID,
				p.Name,
			)
		} else if p.Description != nil && p.Name == nil {
			res, err = s.Store.DB.Exec(
				"UPDATE projects SET description = $2 WHERE id = $1",
				p.ID,
				p.Description,
			)
		} else if p.Description != nil && p.Name != nil {
			res, err = s.Store.DB.Exec(
				"UPDATE projects SET name = $2, description = $3 WHERE id = $1",
				p.ID,
				p.Name,
				p.Description,
			)
		} else {
			err = errors.New("at least one field have to be not blank")
		}
		if err != nil {
			s.errorresp(w, r, http.StatusInternalServerError, err)
			return
		}
		if st, _ := res.RowsAffected(); st == 0 {
			s.errorresp(w, r, http.StatusNotFound, errors.New("project not found"))
			return
		}

		if err := s.Store.DB.QueryRow(
			"SELECT id, name, description FROM projects WHERE id = $1",
			params["id"],
		).Scan(
			&p.ID,
			&p.Name,
			&p.Description,
		); err != nil {
			s.errorresp(w, r, http.StatusInternalServerError, err)
			return
		}

		s.successresp(w, r, http.StatusOK, p)
	}
}

// Delete project
func (s *APIServer) deleteProject() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		params := mux.Vars(r)

		res, err := s.Store.DB.Exec(
			"DELETE FROM projects WHERE id = $1",
			params["id"],
		)

		if err != nil {
			s.errorresp(w, r, http.StatusInternalServerError, err)
			return
		}
		if st, _ := res.RowsAffected(); st == 0 {
			s.errorresp(w, r, http.StatusNotFound, errors.New("project not found"))
			return
		}

		s.successresp(w, r, http.StatusNoContent, nil)
	}
}
