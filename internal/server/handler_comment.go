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

// Add new comment
func (s *APIServer) createComment() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		req := &model.Comment{}
		var err error

		if err = json.NewDecoder(r.Body).Decode(req); err != nil {
			s.errorresp(w, r, http.StatusBadRequest, err)
			return
		}
		if err = req.Validate(); err != nil {
			s.errorresp(w, r, http.StatusBadRequest, err)
			return
		}

		c := &model.Comment{
			ID:     genID(commentID),
			TaskID: req.TaskID,
			Text:   req.Text,
		}

		if err = s.Store.DB.QueryRow(
			"INSERT INTO comments (id, task_id, text) VALUES ($1, $2, $3) RETURNING date",
			c.ID,
			c.TaskID,
			c.Text,
		).Scan(&c.Date); err != nil {
			s.errorresp(w, r, http.StatusUnprocessableEntity, errors.New("task not found"))
			return
		}
		s.successresp(w, r, http.StatusCreated, c)
	}
}

// Get comments list
func (s *APIServer) getComments() http.HandlerFunc {
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

		rows, err := s.Store.DB.Query("SELECT * FROM comments ORDER BY date")
		if err != nil {
			s.errorresp(w, r, http.StatusInternalServerError, err)
			return
		}
		result := []model.Comment{}
		for rows.Next() {
			c := &model.Comment{}
			err = rows.Scan(
				&c.ID,
				&c.TaskID,
				&c.Text,
				&c.Date)
			result = append(result, *c)
		}

		pstart := min(page*pageSize, int(len(result)))
		pfinish := min(page*pageSize+pageSize, int(len(result)))
		s.successresp(w, r, http.StatusOK, map[string]interface{}{
			"total_count": len(result),
			"data":        result[pstart:pfinish],
		})
	}
}

// Get single comment
func (s *APIServer) getComment() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		params := mux.Vars(r)

		p := &model.Comment{}
		if err := s.Store.DB.QueryRow(
			"SELECT id, task_id, text, date FROM comments WHERE id = $1",
			params["id"],
		).Scan(
			&p.ID,
			&p.TaskID,
			&p.Text,
			&p.Date,
		); err != nil {
			s.errorresp(w, r, http.StatusNotFound, errors.New("comment not found"))
			return
		}

		s.successresp(w, r, http.StatusOK, p)
	}
}

// Patch project
func (s *APIServer) patchComment() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		req := &model.Comment{}
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
		c := &model.Comment{
			ID:   &id,
			Text: req.Text,
		}

		var res sql.Result
		if c.Text != nil {
			res, err = s.Store.DB.Exec(
				"UPDATE comments SET text = $2 WHERE id = $1",
				c.ID,
				c.Text,
			)
		} else {
			s.errorresp(w, r,
				http.StatusBadRequest,
				errors.New("text field have to be not blank"),
			)
			return
		}
		if err != nil {
			s.errorresp(w, r, http.StatusInternalServerError, err)
			return
		}
		if st, _ := res.RowsAffected(); st == 0 {
			s.errorresp(
				w, r,
				http.StatusNotFound,
				errors.New("comment not found"),
			)
			return
		}

		if err = s.Store.DB.QueryRow(
			"SELECT id, task_id, text, date FROM comments WHERE id = $1",
			params["id"],
		).Scan(
			&c.ID,
			&c.TaskID,
			&c.Text,
			&c.Date,
		); err != nil {
			s.errorresp(w, r, http.StatusInternalServerError, err)
			return
		}

		s.successresp(w, r, http.StatusOK, c)
	}
}

// Delete project
func (s *APIServer) deleteComment() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		params := mux.Vars(r)

		res, err := s.Store.DB.Exec(
			"DELETE FROM comments WHERE id = $1",
			params["id"],
		)

		if err != nil {
			s.errorresp(w, r, http.StatusInternalServerError, err)
			return
		}
		if st, _ := res.RowsAffected(); st == 0 {
			s.errorresp(w, r, http.StatusNotFound, errors.New("comment not found"))
			return
		}

		s.successresp(w, r, http.StatusNoContent, nil)
	}
}
