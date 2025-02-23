package testenv

// func (e *TestEnvironment) SeedAndLoginUser(t *testing.T) (*domain.User, *tokens.Token) {
// 	t.Helper()

// 	username := cuid.Slug()
// 	user, err := e.Database.Users().CreateUser(context.Background(), domain.CreateUserArgs{
// 		Email:        username + "@letters2me.com",
// 		PasswordHash: []byte("randompassword"),
// 		Username:     username,
// 		IsVerified:   true,
// 	})
// 	require.NoError(t, err)

// 	token, err := tokens.New(tokens.WithTTL(auth.LoginTokenTTL))
// 	require.NoError(t, err)

// 	_, err = e.Database.AuthTokens().CreateToken(context.Background(), domain.CreateTokenArgs{
// 		Hash:      token.Hash,
// 		UserID:    user.ID,
// 		ExpiresAt: token.Expiry.Time,
// 	})
// 	require.NoError(t, err)

// 	return user, token
// }
