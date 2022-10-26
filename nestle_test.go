package nestle_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/tarunKoyalwar/nestle"

	_ "embed"
)

var graphql_query string = `query getNotes(
	$limit: Int
	$offset: Int
	$user_ids: [String!]!
	$note_type: [String!]!
	$exclude_user_ids: [String!]
  ) {
	notes(
	  limit: $limit
	  offset: $offset
	  user_ids: $user_ids
	  note_type: $note_type
	  exclude_user_ids: $exclude_user_ids
	) {
	  data {
		id
		attributes {
		  body
		  note_type
		  published_at
		  resource_id
		  resource_type
		  clients {
			user_canonical_id
			full_name
		  }
		  author {
			user_canonical_id
			full_name
		  }
		  metadata
		}
	  }
	  metaverse {
		total
		offset
		limit
		sort_by
		sort_order
	  }
	}
  }nested
  
  mutation publishNote(
	$body: String!
	$user_ids: [String!]!
	$note_type: String!
  ) {
	publishNote(body: $body, user_ids: $user_ids, note_type: $note_type) {
	  data {
		id
		attributes {
		  body
		  note_type
		  published_at
		  resource_id
		  resource_type
		  clients {
			user_canonical_id
			full_name
		  }
		  author {
			user_canonical_id
			full_name
		  }
		  metadata
		}
	  }
	}
  }nested

`

func Test_SyntaxCheck(t *testing.T) {
	/*
		Prechecks
		- correct prematch and postmatch groups

	*/

	//`(query|mutation) [A-Za-z]+([(]){1}[^({|})]+([)])[\s]*([{])`

	prematch := `(query|mutation)\s+[a-zA-Z]+[0-9]*[a-zA-Z]+(\([^(\(|\))]+\))*\s*`

	nested_regex := prematch + `[{:nested:}]` + `nested`

	n, err := nestle.MustCompile(nested_regex)

	if err != nil {
		t.Fatalf("Failed to compile regex")
	}

	if n.EndDelim != '}' {
		t.Errorf("Failed to correctly parse end delimeter expected } but got %v", n.EndDelim)
	}

	if n.StartDelim != '{' {
		t.Errorf("Failed to correctly parse start delimeter expected { but got %v", n.StartDelim)
	}

	if n.Prematch != prematch {
		t.Errorf("Failed to correctly parse prematch group \nExpected: %v\nGot: %v", prematch, n.Prematch)
	}

}

func Test_NormalCase(t *testing.T) {
	t.Logf("Test case where pre and post match is available")
	prematch := `(query|mutation)\s+[a-zA-Z]+[0-9]*[a-zA-Z]+(\([^(\(|\))]+\))*\s*`
	postmatch := `(nested)`
	nested_regex := prematch + `[{:nested:}]` + postmatch

	re, er := nestle.MustCompile(nested_regex)

	if er != nil {
		t.Fatalf("failed to compile nested regex %v", nested_regex)
	}

	results := re.FindAllString(graphql_query)

	got := strings.Fields(strings.Join(results, "\n"))
	expected := strings.Fields(graphql_query)

	if !reflect.DeepEqual(expected, got) {
		t.Errorf("Nestle failed where both pre and post match were available")
		t.Logf("expected %v\n", expected)
		t.Logf("got %v\n", got)
	} else {
		t.Logf("Test Successfull")
	}

}

func Test_OnlyPrematch(t *testing.T) {
	t.Logf("Test case where only prematch is available")

	nested_regex := `(attributes)\s[{:nested:}]`

	re, er := nestle.MustCompile(nested_regex)

	if er != nil {
		t.Fatalf("failed to compile nested regex %v", nested_regex)
	}
	results := re.FindAllString(graphql_query)

	if len(results) != 2 {
		t.Fatalf("nestle failed expected 2 results but got %v", len(results))
	}

	got := strings.Fields(results[1])
	expected := strings.Fields(`attributes {
        body
        note_type
        published_at
        resource_id
        resource_type
        clients {
          user_canonical_id
          full_name
        }
        author {
          user_canonical_id
          full_name
        }
        metadata
      }
`)

	if !reflect.DeepEqual(expected, got) {
		t.Errorf("Nestle failed where only prematch is available")
		t.Logf("expected %v\n", expected)
		t.Logf("got %v\n", got)
	} else {
		t.Logf("Test Successfull")
	}
}

func Test_OnlyPostmatch(t *testing.T) {
	t.Logf("Test case where only postmatch is available")

	nested_regex := `[{:nested:}]\s*(metaverse)`

	re, er := nestle.MustCompile(nested_regex)

	if er != nil {
		t.Fatalf("failed to compile nested regex %v", nested_regex)
	}
	results := re.FindAllString(graphql_query)

	got := strings.Fields(results[0])
	expected := strings.Fields(`{
		id
		attributes {
		  body
		  note_type
		  published_at
		  resource_id
		  resource_type
		  clients {
			user_canonical_id
			full_name
		  }
		  author {
			user_canonical_id
			full_name
		  }
		  metadata
		}
	  }
	  metaverse 
`)

	if !reflect.DeepEqual(expected, got) {
		t.Errorf("Nestle failed where only postmatch is available")
		t.Logf("expected %v\n", expected)
		t.Logf("got %v\n", got)
	} else {
		t.Logf("Test Successfull")
	}

}
