package apiv1

import (
	"context"
	"reflect"
	"testing"
	"vc/pkg/logger"
	"vc/pkg/model"
	//	"github.com/stretchr/testify/assert"
)

func TestClient_VerifyCredential(t *testing.T) {

	tests := []struct {
		name     string
		fields   fields
		args     args
		want     *VerifyCredentialReply
		errorExp bool
	}{
		newVerifyTestCase("Not a jwt 1", "", false, "not a jwt", false),
		newVerifyTestCase("Not a jwt 2", " ", false, "not a jwt", false),
		newVerifyTestCase("Not a jwt 3", "               .                         .                          ~", false, "not a jwt", false),
		newVerifyTestCase("Not a jwt 4", "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", false, "not a jwt", false),
		newVerifyTestCase("Not a jwt 5", "xxxxxxxxxxxxxxxxxxxxxxxx.xxxxxxxxxxxxxxxxxxxxxxxx.xxxxxxxxxxxxxxxxxxxxxxxx", false, "not a jwt", false),
		newVerifyTestCase("Not a jwt 6", "xxxxxxxxxxxxxxxxxxxxxxxx.xxxxxxxxxxxxxxxxxxxxxxxx.xxxxxxxxxxxxxxxxxxxxxxxx.xxxxxxxxxxxxxxxxxxxxxxxx", false, "not a jwt", false),

		newVerifyTestCase("Header typ does not contains sd-jwt", "eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.xxxxxxxxxxxxxxxxxxxxxxxx.xxxxxxxxxxxxxxxxxxxxxxxx~", false, "supported jwt header.typ are: sd-jwt, vc+sd-jwt", false),

		newVerifyTestCase("Header alg not ES256", "eyJhbGciOiJIUzI1NiIsInR5cCI6InNkLWp3dCJ9.xxxxxxxxxxxxxxxxxxxxxxxx.xxxxxxxxxxxxxxxxxxxxxxxx~", false, "supported jwt header.alg are: ES256", false),

		newVerifyTestCase("Missing or invalid JWK field", "eyJhbGciOiJFUzI1NiIsInR5cCI6InNkLWp3dCJ9.xxxxxxxxxxxxxxxxxxxxxxxx.xxxxxxxxxxxxxxxxxxxxxxxx~", false, "missing or invalid JWK field", false),

		//TODO(mk): add more testcases incl. several successful ones
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
