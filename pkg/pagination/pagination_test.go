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

	t.Run("set result size", func(t *testing.T) {
		p, err := NewSqlPage(nil, nil, 20, 100)
		if err != nil {
			t.Fatal(err)
		}
		p.SetResultSize(5)
		if p.GetPagesCount() != 1 {
			t.Fatal("unexpected total count. expected 1, got", p.GetPagesCount())
		}
	})

	t.Run("set result size for one page", func(t *testing.T) {
		p, err := NewSqlPage(nil, nil, 20, 100)
		if err != nil {
			t.Fatal(err)
		}
		p.SetResultSize(20)
		if p.GetPagesCount() != 1 {
			t.Fatal("unexpected total count. expected 1, got", p.GetPagesCount())
		}

		})

	t.Run("set result size for two pages", func(t *testing.T) {
		p, err := NewSqlPage(nil, nil, 20, 100)
		if err != nil {
			t.Fatal(err)
		}
		p.SetResultSize(21)
		if p.GetPagesCount() != 2 {
			t.Fatal("unexpected total count. expected 2, got", p.GetPagesCount())
		}

		})

	t.Run("get pages count without set size", func(t *testing.T) {
		p, err := NewSqlPage(nil, nil, 1, 100)
		if err != nil {
			t.Fatal(err)
		}
		if p.GetPagesCount() != 0 {
			t.Fatal("unexpected total count")
		}
	})
}

func TestPaginationWithPagesCalculation(t *testing.T) {
	t.Run("total count enough", func(t *testing.T) {
		p, err := NewSqlPage(nil, intPtr(20), 20, 100)
		if err != nil {
			t.Fatal(err)
		}
		p.SetResultSize(40)
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
		p.SetResultSize(40)
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
