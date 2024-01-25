package do

import (
	"testing"
)

func TestSigner_Sign(t *testing.T) {
	type fields struct {
		secret string
	}
	type args struct {
		raw []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantR   string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "1",
			fields: fields{
				secret: func() string {
					return "MIIEvAIBADANBgkqhkiG9w0BAQEFAASCBKYwggSiAgEAAoIBAQC7ZhdgKSxP0nHqLnLRYxg/06aYAta3e/7sEX287uFAq1uvkprgxjU2mxux0x82gvJW6JS8lRWzmYEXeoD8YdKsFAQgf/+vt8jfW918A+jdTYh7306zbfnielUUyYCz3NCRMEnn8j05Ggwj8B4gPYYtxb0ifQldoKsA0InjOSKqq2pT6W2wvvujgqVyVHVqg+Mh4TyT6uSDD65YFD8PVB4xIf4rDrXCVHqdIzQ6KfdaxE4Nv5mCj9ZgM4i2hIbciuBUdb55Doph/y8tu3pNmB9cgAf/3CitkKTFXU20wk4NdyzGlKO0GVW0xwwtXMpY6TRnfitEFIqsDR+gHHwkwg8jAgMBAAECggEAA+aXSaOhw1RB6xomFAYIMlqWpvxheXabHaaBjltklJgx3jY4LhdCHU2zkQswJPjV0DRNgEgUKGW4wliE8yaaytqBztFgl55E5U1VZL2fbRfY5ZyGLVr2hB3VbouBN9KQBC2pTtWD1mVHPRTu85nZlKip4VwRRAHSUr/Z1uvQtRUM5IqAEY+VhzJKa0gQlBUfiRLJVno5xCP60cQ4bHEMM738eQgTw85pmjdR8+n00H6hwhn3CCo9RTkTeeGWkd8Z2kFoHqgS8DWK6/wAtU0Drm5Ql06wM9pyttnpvm9aU7wDFTCTjlTLh2jLZGAunnIYsQ1DldNZwJpxXxo0cPRJeQKBgQDSe2j91Gez/iqHliF5+H2dWoeD68xaPuGqrgSdnre+CluEuplVrOfqBHKDC5g1dWmbJEdGN4nhv8ZofyGjBFmNN2O+AlK0tZ/wqhWmngXe5ac/KQ4llK2tjiojhZCYHKqtAHn2fojxNGc2nFinygQAdwOS9ZEyIX8w4VM1Gg6qiQKBgQDj7MIA8SSAREOT0omlIBKq1AUKZVTideUVoiv+hTJCCkczSEIuqDQoo8bV619r7ObtcxO24pIWYZ7BhLi4+MA/n7IheutziVB9AASnL0OAF2QWQeAGj6W2krERKQO3TOhepSvVKr0tAmzB7IyYqpjsqcVCD1ibc2d/jEqRO0sRSwKBgB7MvEJYcIxgJu0MRP3KJOd/tdDSEmcqSG9nY3mFHjIK5fV4MLPex1jxKaiPa8h20+tF1cAqpFyKaYglAlEOc+Q8NLY7NMsIwMzAtsZY3VcOl/igE1fgd8GryfLEurHnj/oc1bwCLBvPpULSgg6bexZuU/GPSZ3iVPBcKIbet0KxAoGAdtqcZB9beGOglbIhQvFRqrE7G6uxsxHlbv2NUYElrxhq/ov8rxXZdSPKaPz/WmlEFqh+rEzD/1XSknliVlqo7cSaACl4JFyDk1tyEbhsy5vm/lBFwUYhFO6z1Q39ORWqysf47oUF9zWffxSaUHYNnsP88DDOdOmeG/4NWGSCBbcCgYB67BNdEzMiO24gQW/rcIsSLyBBJXzALSgYOqWEWUFocnDXZb6uG0rW8a40sewxahOyLnSNbCCBlQoRQb5oag31M4E+yxE6AjEb4btZcbEHwmjZ+nYGGd56nZrKF5dtGW1iMAQy6YBqaQ/BnQJfpJxi8R4i5E5ZgHn8IebNR23MWg=="
				}(),
			},
			args: args{
				raw: []byte("hello"),
			},
			wantR:   "E2EdXjzaxnDewsPXRlO0L6bdNCUaQkxWG0w1IZIfs9T7w4+6Z5LUJbJ7WDHIVR3RomI0rW7Vv3vRWHB3WkZscpF5EL8o62vBitl5iuu/EUb9a8obWAtxq5Hum4giVZa+15i9dx2k8JN8JuUH1Ug+9+6z+mTM/9a7LrtPR7Gq1MK6ejqoQmnmzKoWaBEJQJeRH/9MlwxW0hc1V9JBaDhmXEJR8d1h53RBLoNUFBLjA2Aaz6zUrBkjZTRN09q8CEjOYWJ1uMX31VdSiiuwYr3novc4EfoJjfKUs7KYm00riDqpi+Sjlqk6vs0Ur6SKiQTPYVfhO3gw10a+FZDvzUlVtQ==",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Signer{
				secret: tt.fields.secret,
			}
			gotR, err := s.Sign(tt.args.raw)
			if (err != nil) != tt.wantErr {
				t.Errorf("Signer.Sign() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotR != tt.wantR {
				t.Errorf("Signer.Sign() = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}
