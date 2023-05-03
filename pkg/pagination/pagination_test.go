package pagination

import "testing"

func TestNewSqlPage(t *testing.T) {
	t.Run("token", func(t *testing.T) {
		p, err := NewSqlPage(nil, nil, 1, 100)
		if err != nil {
			t.Fatal(err)
		}
		if p.Token() != "1" {
			t.Fatal("unexpected token")
		}
	})

	t.Run("size", func(t *testing.T) {
		p, err := NewSqlPage(nil, nil, 1, 100)
		if err != nil {
			t.Fatal(err)
		}
		if p.Size() != 1 {
			t.Fatal("unexpected size")
		}
	})

	t.Run("offset", func(t *testing.T) {
		p, err := NewSqlPage(nil, nil, 1, 100)
		if err != nil {
			t.Fatal(err)
		}
		if p.Offset() != 0 {
			t.Fatal("unexpected offset")
		}
	})

	t.Run("next token", func(t *testing.T) {
		p, err := NewSqlPage(nil, nil, 2, 100)
		if err != nil {
			t.Fatal(err)
		}
		nextToken := p.NextToken(2)
		if nextToken == nil {
			t.Fatal("unexpected next token")
		}
		if *nextToken != "2" {
			t.Fatal("unexpected next token")
		}
	})

	t.Run("next token with current token", func(t *testing.T) {
		p, err := NewSqlPage(strPtr("2"), nil, 2, 100)
		if err != nil {
			t.Fatal(err)
		}
		if p.NextToken(2) == nil {
			t.Fatal("unexpected next token")
		}
	})

	t.Run("invalid token", func(t *testing.T) {
		_, err := NewSqlPage(strPtr("a"), nil, 1, 100)
		if err == nil {
			t.Fatal("unexpected error")
		}
	})

	t.Run("negative token", func(t *testing.T) {
		_, err := NewSqlPage(strPtr("-1"), nil, 1, 100)
		if err == nil {
			t.Fatal("unexpected error")
		}
	})

	t.Run("negative size", func(t *testing.T) {
		_, err := NewSqlPage(nil, intPtr(-1), 1, 100)
		if err == nil {
			t.Fatal("unexpected error")
		}
	})

	t.Run("size greater than max size", func(t *testing.T) {
		_, err := NewSqlPage(nil, intPtr(101), 1, 100)
		if err != nil {
			t.Fatal("unexpected error")
		}
	})

	t.Run("result size less than size", func(t *testing.T) {
		p, err := NewSqlPage(nil, intPtr(10), 10, 100)
		if err != nil {
			t.Fatal("unexpected error")
		}
		if p.NextToken(2) != nil {
			t.Fatal("unexpected next token")
		}
	})

	t.Run("set total count", func(t *testing.T) {
		p, err := NewSqlPage(nil, nil, 1, 100)
		if err != nil {
			t.Fatal(err)
		}
		p.SetTotalCount(1)
		if p.GetTotalCount() != 1 {
			t.Fatal("unexpected total count")
		}
	})

	t.Run("get total count", func(t *testing.T) {
		p, err := NewSqlPage(nil, nil, 1, 100)
		if err != nil {
			t.Fatal(err)
		}
		if p.GetTotalCount() != 0 {
			t.Fatal("unexpected total count")
		}
	})
}

func TestPaginationWithTotalCount(t *testing.T) {
	t.Run("total count enough", func(t *testing.T) {
		p, err := NewSqlPage(nil, intPtr(20), 20, 100)
		if err != nil {
			t.Fatal(err)
		}
		p.SetTotalCount(40)
		nextToken := p.NextToken(20)
		if nextToken == nil {
			t.Fatal("unexpected next token")
		}
		if *nextToken != "2" {
			t.Fatal("unexpected next token")
		}
	})

	t.Run("total count not enough", func(t *testing.T) {
		p, err := NewSqlPage(strPtr("2"), intPtr(20), 20, 100)
		if err != nil {
			t.Fatal(err)
		}
		p.SetTotalCount(40)
		nextToken := p.NextToken(20)
		if nextToken != nil {
			t.Fatal("unexpected next token")
		}
	})
}

func intPtr(i int) *int {
	return &i
}

func strPtr(s string) *string {
	return &s
}
