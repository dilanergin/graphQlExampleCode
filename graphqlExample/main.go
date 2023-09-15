package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/graphql-go/graphql"
)

type Book struct {
	ID        int
	Title     string
	Author    Author
	Year      int
	PageCount int
	Comments  []Comment
}

type Author struct {
	Name  string
	Books []int
}

type Comment struct {
	Body string
}

func bookData() []Book {
	author := &Author{Name: "Ahmet Ümit", Books: []int{1}}
	book1 := Book{
		ID:     1,
		Title:  "Sis ve Gece",
		Author: *author,
		Comments: []Comment{
			Comment{Body: "Güzel bir kitap"},
			Comment{Body: "sevdim"}},
		PageCount: 200,
		Year:      2012,
	}
	book2 := Book{
		ID:        2,
		Title:     "Aşk Masalı",
		Author:    *author,
		Comments:  []Comment{Comment{Body: "Ortalama"}},
		PageCount: 289,
		Year:      2022,
	}

	var books []Book
	books = append(books, book1, book2)

	return books
}

func main() {
	books := bookData()
	var commentType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Comment",

			Fields: graphql.Fields{
				"body": &graphql.Field{
					Type: graphql.String,
				},
			},
		},
	)
	var authorType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Author",
			Fields: graphql.Fields{
				"Name": &graphql.Field{
					Type: graphql.String,
				},
				"Books": &graphql.Field{
					Type: graphql.NewList(graphql.Int),
				},
			},
		},
	)
	var bookType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Book",
			Fields: graphql.Fields{
				"id":        &graphql.Field{Type: graphql.Int},
				"title":     &graphql.Field{Type: graphql.String},
				"author":    &graphql.Field{Type: authorType},
				"comments":  &graphql.Field{Type: graphql.NewList(commentType)},
				"pageCount": &graphql.Field{Type: graphql.Int},
				"year":      &graphql.Field{Type: graphql.Int},
			},
		},
	)

	fields := graphql.Fields{
		"book": &graphql.Field{
			Type:        bookType,
			Description: "Get Book By ID",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id, ok := p.Args["id"].(int)
				if ok {
					for _, book := range books {
						if int(book.ID) == id {

							return book, nil
						}
					}
				}
				return nil, nil
			},
		},

		"list": &graphql.Field{
			Type:        graphql.NewList(bookType),
			Description: "Get Book List",
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				return books, nil
			},
		},
	}
	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatalf("Errors: %v", err)
	}
	query := `
        {
			book(id:1) {
				title
				author {
					Name
					Books
				}
				pageCount
			}
        }
    `
	params := graphql.Params{Schema: schema, RequestString: query}
	r := graphql.Do(params)
	if len(r.Errors) > 0 {
		log.Fatalf("Errors: %+v", r.Errors)
	}
	rJSON, _ := json.Marshal(r)
	fmt.Printf("%s \n", rJSON) // {“data”:{“hello”:”world”}}
}
