/*
 * People info
 *
 * No description provided (generated by Swagger Codegen https://github.com/swagger-api/swagger-codegen)
 *
 * API version: 0.0.1
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package swagger

type People struct {
	Surname string `json:"surname"`
	Name string `json:"name"`
	Patronymic string `json:"patronymic,omitempty"`
	Address string `json:"address"`
}
