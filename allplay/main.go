package main

import (
	"bytes"
	"io"
	"log"
	"net/http"

	"github.com/hibooboo2/ggames/allplay/pollen"
)

func main() {
	log.SetFlags(log.Lshortfile)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		g := pollen.NewGame([]pollen.PlayerInput{
			{"JAMES", nil},
			{"RAE", nil},
		})

		err := g.PlayCard("JAMES", g.GetHand("JAMES")[0].ID, pollen.Position{0.5, 0.5})
		if err != nil {
			panic(err)
		}
		err = g.NextPlayer()
		if err != nil {
			panic(err)
		}

		err = g.PlayCard("RAE", g.GetHand("RAE")[0].ID, pollen.Position{0.5, -0.5})
		if err != nil {
			panic(err)
		}

		tk := g.GetNextToken()
		err = g.PlayToken("RAE", tk, pollen.Position{1, 0})
		if err != nil {
			panic(err)
		}

		err = g.NextPlayer()
		if err != nil {
			panic(err)
		}

		err = g.PlayCard("JAMES", g.GetHand("JAMES")[0].ID, pollen.Position{-0.5, 0.5})
		if err != nil {
			panic(err)
		}

		tk = g.GetNextToken()
		err = g.PlayToken("JAMES", tk, pollen.Position{0, 1})
		if err != nil {
			panic(err)
		}

		err = g.NextPlayer()
		if err != nil {
			panic(err)
		}

		err = g.PlayCard("RAE", g.GetHand("RAE")[0].ID, pollen.Position{-0.5, -0.5})
		if err != nil {
			panic(err)
		}

		tk = g.GetNextToken()
		err = g.PlayToken("RAE", tk, pollen.Position{0, -1})
		if err != nil {
			panic(err)
		}

		tk = g.GetNextToken()
		err = g.PlayToken("RAE", tk, pollen.Position{-1, 0})
		if err != nil {
			panic(err)
		}

		err = g.NextPlayer()
		if err != nil {
			panic(err)
		}

		err = g.PlayCard("JAMES", g.GetHand("JAMES")[0].ID, pollen.Position{-1.5, -0.5})
		if err != nil {
			panic(err)
		}

		tk = g.GetNextToken()
		err = g.PlayToken("JAMES", tk, pollen.Position{-1, -1})
		if err != nil {
			panic(err)
		}

		err = g.NextPlayer()
		if err != nil {
			panic(err)
		}

		err = g.PlayCard("RAE", g.GetHand("RAE")[0].ID, pollen.Position{-1.5, -1.5})
		if err != nil {
			panic(err)
		}

		tk = g.GetNextToken()
		err = g.PlayToken("RAE", tk, pollen.Position{-2, -1})
		if err != nil {
			panic(err)
		}

		err = g.NextPlayer()
		if err != nil {
			panic(err)
		}

		// err = g.PlayCard("JAMES", g.GetHand("JAMES")[0].ID, pollen.Position{-1.5, -1.5})
		// if err != nil {
		// 	panic(err)
		// }

		buff := bytes.NewBuffer(nil)

		wr := io.MultiWriter(buff, w)
		err = g.Render(wr)
		if err != nil {
			panic(err)
		}
		// log.Println(buff.String())
	})

	// http.Handle("/images", http.StripPrefix("./pollen/images", http.FileServer(http.Dir("./pollen/images"))))
	http.HandleFunc("/images/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./pollen/images/card.png")
	})

	http.ListenAndServe(":8080", nil)
}
