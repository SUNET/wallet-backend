package apiv1

import (
	"context"
	"reflect"
	"testing"
	"vc/pkg/logger"
	"vc/pkg/model"
	//	"github.com/stretchr/testify/assert"
)

// TODO(mk): generate a new version of VALID_EHIC_SD_JWT after bugs have been fixed and merged to main

// VALID_EHIC_SD_JWT with three disclosures and exp 100 years from 20241113
//
// ATTENTION: created using mockas while bugs with attibute-names existed, for example: givenName in disclosure and firstName in claim...
//
// "alg": "ES256",
// "typ": "sd-jwt"
// "_sd_alg": "sha-256"
//
//	"cnf": {
//	  "jwk": {
//	    "crv": "P-256",
//	    "d": "V5rwAvUIx_it5yA9CTZVAsca14b9kjAibM3GsoKn_wk",
//	    "kid": "singing_",
//	    "kty": "EC",
//	    "x": "Cdj7z0qgkhiDqUcdQLPH4c3h3icT4bOP5aIjjU_hu-I",
//	    "y": "Fy_vz3wnc5t1oWIMXtaShexUddY-EynROdaiDRtdf60"
//	  }
//	}
const VALID_EHIC_SD_JWT = "eyJhbGciOiJFUzI1NiIsInR5cCI6InNkLWp3dCJ9.eyJfc2RfYWxnIjoic2hhLTI1NiIsImNhcmRIb2xkZXIiOnsiX3NkIjpbIk5tUmhZekkwTlRZellqa3lZelpqT0dOaU1qRTFOR0l6TUdSa05HRTBNV05pWWpZME5EWXlOemhoWmpOall6bG1Oak0wTVRFd1pqSm1aV1ZoTkRCak5RIiwiWkdJeVkyRmtZMlUwTXpnMVlUazRPRGswT1RBMllUY3lNRFkzWXpjME1HWmhOamxtTXpZMU5ETTJaV0ptTjJKbE9HSXhPR1EwTkRjNU5tRTJNemRqTVEiLCJNbVEyTm1SbU5UUmpNams0WldGbE0yWTNNelkxT1dKbVpqSTVaakV4TWpVM05UUXhaakprT1RFMk5qUTNOVEEzTmpBNFpHTTJPV05pTVdWaVpqQmlOUSJdLCJjYXJkaG9sZGVyU3RhdHVzIjoiYWN0aXZlIiwiaWQiOiJjYjQ3MGE2My0yYTQxLTQ1YzUtYmMyMC0xYmQyMDlkMGI0ZjUifSwiY2FyZEluZm9ybWF0aW9uIjp7ImV4cGlyeURhdGUiOiIxOTY5LTA4LTE3IDE3OjMxOjQxLjM2MjkzNzg4NCArMDAwMCBVVEMiLCJpZCI6ImM2NmU3MmU0LTgwMTItNDUxNC04ZjA2LTBmNTNlZWZjNDFhZCIsImludmFsaWRTaW5jZSI6IjIwMDEtMDYtMjggMTc6MzM6MTMuMjkxODMyNjYgKzAwMDAgVVRDIiwiaXNzdWFuY2VEYXRlIjoiMTkxMS0wNS0xMSAwNzo1MTozMy4xODY1MTUzMDcgKzAwMDAgVVRDIiwic2lnbmF0dXJlIjp7Imlzc3VlciI6Ikp1bnlvIiwic2VhbCI6ImZhNzBhMDI3LWU1YjItNDgzZS1iNmQ1LWRmYjdlZWMzYjJiYyJ9LCJ2YWxpZFNpbmNlIjoiMjAxNi0wOC0xOCAxMTo1MDozOS40NTM3NzUxNjcgKzAwMDAgVVRDIn0sImNuZiI6eyJqd2siOnsiY3J2IjoiUC0yNTYiLCJkIjoiVjVyd0F2VUl4X2l0NXlBOUNUWlZBc2NhMTRiOWtqQWliTTNHc29Lbl93ayIsImtpZCI6InNpbmdpbmdfIiwia3R5IjoiRUMiLCJ4IjoiQ2RqN3owcWdraGlEcVVjZFFMUEg0YzNoM2ljVDRiT1A1YUlqalVfaHUtSSIsInkiOiJGeV92ejN3bmM1dDFvV0lNWHRhU2hleFVkZFktRXluUk9kYWlEUnRkZjYwIn19LCJjb21wZXRlbnRJbnN0aXR1dGlvbiI6eyJpZCI6Ijk5NDQ2OWFlLTAwM2MtNGUyZi04MTRkLWVhZmMzOGRjZTdmMiIsImluc3RpdHV0aW9uTmFtZSI6IlhhdG9yaSJ9LCJleHAiOjQ4ODUwOTg1OTgsImlzcyI6Imh0dHBzOi8vaXNzdWVyLnN1bmV0LnNlIiwibmJmIjoxNzMxNDk4NTk4LCJwaWQiOnsiZXhoaWJpdG9ySUQiOiI1ODE4Mjg1NzMwIiwiZmlyc3ROYW1lIjoiQmV0c3kiLCJnZW5kZXIiOiJmZW1hbGUiLCJsYXN0TmFtZSI6IkR1cmdhbiJ9LCJzaWduYXR1cmUiOnsiaXNzdWVyIjoiTmllbHNlbiIsInNlYWwiOiIxOWEyNjllYi1jZjM0LTQ3ZTctODhkYy1jNTE5YmZlN2FiOGEifSwic3RhdHVzIjoiIiwidmN0IjoiaHR0cHM6Ly9jcmVkZW50aWFsLnN1bmV0LnNlL2lkZW50aXR5X2NyZWRlbnRpYWwifQ.0fp3XCnGhLRO8WeU0unjG4yMGID0FdvQhcGTUYCuj-MkP5xgX_IdEASsm4NbGtYJVw-lMMOPGEajKGvz67SMNQ~WyJNVDFPUDBCOEpUVXBPU2RmS1Qxb1RXbyIsImdpdmVuTmFtZSIsIkJldHN5Il0~WyJWbFJ1VVNwa1JWZHFQMVpUYkd0OVdHUSIsImJpcnRoRGF0ZSIsIjE5NjEtMDItMjggMDE6MTY6NDQuNjgwODY4OTI3ICswMDAwIFVUQyJd~WyJjQzg3V2owcE5FVXNPVEJaTlVna1hETSIsImZhbWlseU5hbWUiLCJEdXJnYW4iXQ~"

// TAMPERED_PAYLOAD_EHIC_SD_JWT same as VALID_EHIC_SD_JWT but one char in jwt payload has been changed
const TAMPERED_PAYLOAD_EHIC_SD_JWT = "eyJhbGciOiJFUzI1NiIsInR5cCI6InNkLWp3dCJ9.eyJfc2RfYWxxIjoic2hhLTI1NiIsImNhcmRIb2xkZXIiOnsiX3NkIjpbIk5tUmhZekkwTlRZellqa3lZelpqT0dOaU1qRTFOR0l6TUdSa05HRTBNV05pWWpZME5EWXlOemhoWmpOall6bG1Oak0wTVRFd1pqSm1aV1ZoTkRCak5RIiwiWkdJeVkyRmtZMlUwTXpnMVlUazRPRGswT1RBMllUY3lNRFkzWXpjME1HWmhOamxtTXpZMU5ETTJaV0ptTjJKbE9HSXhPR1EwTkRjNU5tRTJNemRqTVEiLCJNbVEyTm1SbU5UUmpNams0WldGbE0yWTNNelkxT1dKbVpqSTVaakV4TWpVM05UUXhaakprT1RFMk5qUTNOVEEzTmpBNFpHTTJPV05pTVdWaVpqQmlOUSJdLCJjYXJkaG9sZGVyU3RhdHVzIjoiYWN0aXZlIiwiaWQiOiJjYjQ3MGE2My0yYTQxLTQ1YzUtYmMyMC0xYmQyMDlkMGI0ZjUifSwiY2FyZEluZm9ybWF0aW9uIjp7ImV4cGlyeURhdGUiOiIxOTY5LTA4LTE3IDE3OjMxOjQxLjM2MjkzNzg4NCArMDAwMCBVVEMiLCJpZCI6ImM2NmU3MmU0LTgwMTItNDUxNC04ZjA2LTBmNTNlZWZjNDFhZCIsImludmFsaWRTaW5jZSI6IjIwMDEtMDYtMjggMTc6MzM6MTMuMjkxODMyNjYgKzAwMDAgVVRDIiwiaXNzdWFuY2VEYXRlIjoiMTkxMS0wNS0xMSAwNzo1MTozMy4xODY1MTUzMDcgKzAwMDAgVVRDIiwic2lnbmF0dXJlIjp7Imlzc3VlciI6Ikp1bnlvIiwic2VhbCI6ImZhNzBhMDI3LWU1YjItNDgzZS1iNmQ1LWRmYjdlZWMzYjJiYyJ9LCJ2YWxpZFNpbmNlIjoiMjAxNi0wOC0xOCAxMTo1MDozOS40NTM3NzUxNjcgKzAwMDAgVVRDIn0sImNuZiI6eyJqd2siOnsiY3J2IjoiUC0yNTYiLCJkIjoiVjVyd0F2VUl4X2l0NXlBOUNUWlZBc2NhMTRiOWtqQWliTTNHc29Lbl93ayIsImtpZCI6InNpbmdpbmdfIiwia3R5IjoiRUMiLCJ4IjoiQ2RqN3owcWdraGlEcVVjZFFMUEg0YzNoM2ljVDRiT1A1YUlqalVfaHUtSSIsInkiOiJGeV92ejN3bmM1dDFvV0lNWHRhU2hleFVkZFktRXluUk9kYWlEUnRkZjYwIn19LCJjb21wZXRlbnRJbnN0aXR1dGlvbiI6eyJpZCI6Ijk5NDQ2OWFlLTAwM2MtNGUyZi04MTRkLWVhZmMzOGRjZTdmMiIsImluc3RpdHV0aW9uTmFtZSI6IlhhdG9yaSJ9LCJleHAiOjQ4ODUwOTg1OTgsImlzcyI6Imh0dHBzOi8vaXNzdWVyLnN1bmV0LnNlIiwibmJmIjoxNzMxNDk4NTk4LCJwaWQiOnsiZXhoaWJpdG9ySUQiOiI1ODE4Mjg1NzMwIiwiZmlyc3ROYW1lIjoiQmV0c3kiLCJnZW5kZXIiOiJmZW1hbGUiLCJsYXN0TmFtZSI6IkR1cmdhbiJ9LCJzaWduYXR1cmUiOnsiaXNzdWVyIjoiTmllbHNlbiIsInNlYWwiOiIxOWEyNjllYi1jZjM0LTQ3ZTctODhkYy1jNTE5YmZlN2FiOGEifSwic3RhdHVzIjoiIiwidmN0IjoiaHR0cHM6Ly9jcmVkZW50aWFsLnN1bmV0LnNlL2lkZW50aXR5X2NyZWRlbnRpYWwifQ.0fp3XCnGhLRO8WeU0unjG4yMGID0FdvQhcGTUYCuj-MkP5xgX_IdEASsm4NbGtYJVw-lMMOPGEajKGvz67SMNQ~WyJNVDFPUDBCOEpUVXBPU2RmS1Qxb1RXbyIsImdpdmVuTmFtZSIsIkJldHN5Il0~WyJWbFJ1VVNwa1JWZHFQMVpUYkd0OVdHUSIsImJpcnRoRGF0ZSIsIjE5NjEtMDItMjggMDE6MTY6NDQuNjgwODY4OTI3ICswMDAwIFVUQyJd~WyJjQzg3V2owcE5FVXNPVEJaTlVna1hETSIsImZhbWlseU5hbWUiLCJEdXJnYW4iXQ~"

func TestClient_VerifyCredential(t *testing.T) {
	tests := []struct {
		name     string
		fields   fields
		args     args
		want     *VerifyCredentialReply
		errorExp bool
	}{
		newVerifyTestCase("Not credential provided 1", "", false, MsgNoCredentialProvided, false),
		newVerifyTestCase("Not credential provided 2", " ", false, MsgNoCredentialProvided, false),

		newVerifyTestCase("Not a jwt 1", "               .                         .                          ~", false, MsgNotAJwt, false),
		newVerifyTestCase("Not a jwt 2", "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", false, MsgNotAJwt, false),
		newVerifyTestCase("Not a jwt 3", "xxxxxxxxxxxxxxxxxxxxxxxx.xxxxxxxxxxxxxxxxxxxxxxxx.xxxxxxxxxxxxxxxxxxxxxxxx", false, MsgNotAJwt, false),
		newVerifyTestCase("Not a jwt 4", "xxxxxxxxxxxxxxxxxxxxxxxx.xxxxxxxxxxxxxxxxxxxxxxxx.xxxxxxxxxxxxxxxxxxxxxxxx.xxxxxxxxxxxxxxxxxxxxxxxx", false, MsgNotAJwt, false),
		newVerifyTestCase("Not a jwt 5", "xxxxxxxxxxxxxxxxxxxxxxxx.xxxxxxxxxxxxxxxxxxxxxxxx_xxxxxxxxxxxxxxxxxxxxxxxx", false, MsgNotAJwt, false),

		newVerifyTestCase("Header typ does not contains sd-jwt", "eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.xxxxxxxxxxxxxxxxxxxxxxxx.xxxxxxxxxxxxxxxxxxxxxxxx~", false, "supported jwt header.typ are: sd-jwt, vc+sd-jwt", false),

		newVerifyTestCase("Header alg not ES256", "eyJhbGciOiJIUzI1NiIsInR5cCI6InNkLWp3dCJ9.xxxxxxxxxxxxxxxxxxxxxxxx.xxxxxxxxxxxxxxxxxxxxxxxx~", false, "supported jwt header.alg are: ES256", false),

		newVerifyTestCase("Missing or invalid JWK", "eyJhbGciOiJFUzI1NiIsInR5cCI6InNkLWp3dCJ9.xxxxxxxxxxxxxxxxxxxxxxxx.xxxxxxxxxxxxxxxxxxxxxxxx~", false, MsgInvalidJwk, false),

		newVerifyTestCase("Tampered jwt payload in EHIC", TAMPERED_PAYLOAD_EHIC_SD_JWT, false, MsgUnableToParseToken, false),

		newVerifyTestCase("Valid EHIC", VALID_EHIC_SD_JWT, true, "", false),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				cfg: tt.fields.cfg,
				log: tt.fields.log,
			}

			got, err := c.VerifyCredential(tt.args.ctx, tt.args.request)

			if (err != nil) && tt.errorExp {
				return
			} else if err != nil {
				t.Errorf("VerifyCredential() error is expected to return error")
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("VerifyCredential() got = %v, want %v", got, tt.want)
			}
		})
	}
}

type fields struct {
	cfg *model.Cfg
	log *logger.Log
}

type args struct {
	ctx     context.Context
	request *Credential
}

func newVerifyTestCase(name, credential string, valid bool, message string, errorExp bool) struct {
	name     string
	fields   fields
	args     args
	want     *VerifyCredentialReply
	errorExp bool
} {
	return struct {
		name     string
		fields   fields
		args     args
		want     *VerifyCredentialReply
		errorExp bool
	}{
		name:   name,
		fields: defaultFields(),
		args: args{
			ctx:     context.Background(),
			request: &Credential{Credential: credential},
		},
		want: &VerifyCredentialReply{
			Valid:   valid,
			Message: message,
		},
		errorExp: errorExp,
	}
}

func defaultFields() fields {
	return fields{
		cfg: &model.Cfg{},
		log: logger.NewSimple("testing_handler"),
	}
}
