package authorization_test

// type MockIdentity struct {
// 	ID uint
// }

// // GetToken implements authorization_iface.Identity.
// func (m *MockIdentity) GetToken(appname string, secret string) (string, error) {
// 	panic("unimplemented")
// }

// // GetAgentType implements authorization_iface.Identity.
// func (m *MockIdentity) GetAgentType() identity_iface.AgentType {
// 	return identity_iface.TestAgent
// }

// // IsTokenExpired implements authorization_iface.Identity.
// func (m *MockIdentity) IsTokenExpired(tx *gorm.DB) (bool, error) {
// 	return false, nil
// }

// // GetExpired implements authorization_iface.Identity.
// func (m *MockIdentity) GetExpired(tx *gorm.DB) (*authorization_iface.ExpiredToken, error) {
// 	return &authorization_iface.ExpiredToken{
// 		LastCreated: 1,
// 		LastReset:   0,
// 	}, nil
// }

// // GetUserID implements authorization_iface.Identity.
// func (m *MockIdentity) GetUserID() uint {
// 	return m.ID
// }

// // IdentityID implements authorization_iface.Identity.
// func (m *MockIdentity) IdentityID() uint {
// 	return m.ID
// }

// // IsSuperUser implements authorization_iface.Identity.
// func (m *MockIdentity) IsSuperUser() bool {
// 	return false
// }

// // IsSuperUser implements authorization_iface.Identity.
// func (u *MockIdentity) HasRole(db *gorm.DB, domainID uint, keyname string) (bool, error) {
// 	return true, nil
// }

// type MockOrder struct{}

// // GetEntityID implements authorization.Entity.
// func (m *MockOrder) GetEntityID() string {
// 	return "order"
// }

// func TestAuhtorization(t *testing.T) {

// 	cs := &MockIdentity{}

// 	database_mock.NewDBScenario(t, func(db *database_mock.DBScen) {
// 		auth := authorization.NewAuthorization(ware_cache.NewLocalCache(), db.Tx, "test_phrase")

// 		t.Run("test auth tanpa apapun", func(t *testing.T) {
// 			cache, err := auth.CheckPermission(cs, authorization_iface.CheckPermissionGroup{
// 				&MockOrder{}: &authorization_iface.CheckPermission{
// 					DomainID: 1,
// 					Actions: []authorization_iface.Action{
// 						authorization_iface.Create,
// 					},
// 				},
// 			})

// 			assert.NotNil(t, err)
// 			assert.False(t, cache)
// 		})

// 		t.Run("test di domain", func(t *testing.T) {

// 			domain := auth.Domain(2)
// 			t.Run("test create role", func(t *testing.T) {

// 				role, err := domain.CreateRole("cs")
// 				assert.Nil(t, err)
// 				assert.NotEmpty(t, role)
// 				assert.NotEqual(t, 0, role.ID)

// 				// t.Run("test double create role", func(t *testing.T) {
// 				// 	role, err := domain.CreateRole("cs")
// 				// 	assert.NotEmpty(t, err)
// 				// 	assert.NotEmpty(t, role)
// 				// })

// 				t.Run("test list role", func(t *testing.T) {
// 					roles, err := domain.ListRole()
// 					assert.Nil(t, err)
// 					assert.NotEmpty(t, roles)
// 					assert.Len(t, roles, 1)

// 				})

// 				t.Run("test add role permission", func(t *testing.T) {
// 					err := domain.RoleAddPermission(role.ID, authorization.AddPermissionPayload{
// 						&MockOrder{}: []*authorization.AddPermItem{
// 							{
// 								Action: authorization_iface.Create,
// 								Policy: authorization_iface.Allow,
// 							},
// 							{
// 								Action: authorization_iface.Delete,
// 								Policy: authorization_iface.Allow,
// 							},
// 						},
// 					})
// 					assert.Nil(t, err)
// 				})

// 				t.Run("test role list permission", func(t *testing.T) {
// 					list, err := domain.RoleListPermission(role.ID)
// 					assert.Nil(t, err)
// 					assert.NotEmpty(t, list)
// 					assert.Len(t, list, 2)

// 				})

// 				t.Run("test role remove permission", func(t *testing.T) {
// 					err := domain.RoleRemovePermission(role.ID, authorization.AddPermissionPayload{
// 						&MockOrder{}: []*authorization.AddPermItem{
// 							{
// 								Action: authorization_iface.Delete,
// 								Policy: authorization_iface.Allow,
// 							},
// 						},
// 					})
// 					assert.Nil(t, err)

// 					list, err := domain.RoleListPermission(role.ID)
// 					assert.Nil(t, err)
// 					assert.NotEmpty(t, list)
// 					assert.Len(t, list, 1)
// 				})

// 				t.Run("test user add role", func(t *testing.T) {
// 					err := domain.UserAddRole(cs, role.ID)
// 					assert.Nil(t, err)

// 				})

// 				t.Run("test check user have access", func(t *testing.T) {
// 					cache, err := auth.CheckPermission(cs, authorization_iface.CheckPermissionGroup{
// 						&MockOrder{}: &authorization_iface.CheckPermission{
// 							DomainID: domain.ID,
// 							Actions: []authorization_iface.Action{
// 								authorization_iface.Create,
// 							},
// 						},
// 					})

// 					assert.Nil(t, err)
// 					assert.False(t, cache)

// 				})

// 				t.Run("user list role", func(t *testing.T) {
// 					list, err := domain.UserListRole(cs.IdentityID())
// 					assert.Nil(t, err)
// 					assert.NotEmpty(t, list)
// 				})

// 				t.Run("test user remove role", func(t *testing.T) {
// 					err := domain.UserRemoveRole(cs.IdentityID(), role.ID)
// 					assert.Nil(t, err)
// 				})

// 				t.Run("test check user not have access", func(t *testing.T) {
// 					auth.FlushCache()
// 					cache, err := auth.CheckPermission(cs, authorization_iface.CheckPermissionGroup{
// 						&MockOrder{}: &authorization_iface.CheckPermission{
// 							DomainID: domain.ID,
// 							Actions: []authorization_iface.Action{
// 								authorization_iface.Create,
// 							},
// 						},
// 					})

// 					assert.NotNil(t, err)
// 					assert.False(t, cache)
// 				})

// 				t.Run("test delete role", func(t *testing.T) {
// 					err := domain.DeleteRole(role.ID)
// 					assert.Nil(t, err)
// 				})

// 			})
// 		})
// 	})

// }
