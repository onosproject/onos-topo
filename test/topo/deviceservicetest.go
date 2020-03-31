// Copyright 2019-present Open Networking Foundation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package topo

import (
	"context"
	"crypto/tls"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io"
	"testing"

	"github.com/onosproject/onos-topo/api/device"
	"github.com/stretchr/testify/assert"
)

// TestDeviceService : test
func (s *TestSuite) TestDeviceService(t *testing.T) {
	client, err := getDeviceServiceClient()
	assert.NoError(t, err)

	list, err := client.List(context.Background(), &device.ListRequest{})
	assert.NoError(t, err)

	count := 0
	for {
		_, err := list.Recv()
		if err == io.EOF {
			break
		}
		assert.NoError(t, err)
		count++
	}

	assert.Equal(t, 0, count)

	events := make(chan *device.ListResponse)
	eventList, eventErr := client.List(context.Background(), &device.ListRequest{
		Subscribe: true,
	})
	assert.NoError(t, eventErr)

	go func() {
		for {
			response, err := eventList.Recv()
			if err != nil {
				break
			}
			events <- response
		}
	}()

	addResponse, err := client.Add(context.Background(), &device.AddRequest{
		Device: &device.Device{
			ID:      "test1",
			Type:    "Stratum",
			Address: "device-test1:5000",
			Target:  "device-test1",
			Version: "1.0.0",
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, device.ID("test1"), addResponse.Device.ID)
	assert.NotEqual(t, device.Revision(0), addResponse.Device.Revision)

	getResponse, err := client.Get(context.Background(), &device.GetRequest{
		ID: "test1",
	})
	assert.NoError(t, err)

	assert.Equal(t, device.ID("test1"), getResponse.Device.ID)
	assert.Equal(t, addResponse.Device.Revision, getResponse.Device.Revision)

	eventResponse := <-events

	deviceEventTypes := []device.ListResponse_Type{device.ListResponse_NONE, device.ListResponse_ADDED}
	assert.Contains(t, deviceEventTypes, eventResponse.Type)
	assert.Equal(t, device.ID("test1"), eventResponse.Device.ID)
	assert.Equal(t, addResponse.Device.Revision, eventResponse.Device.Revision)

	list, err = client.List(context.Background(), &device.ListRequest{})
	assert.NoError(t, err)
	for {
		response, err := list.Recv()
		if err == io.EOF {
			break
		}
		assert.NoError(t, err)
		assert.Equal(t, device.ListResponse_NONE, response.Type)
		assert.Equal(t, device.ID("test1"), response.Device.ID)
		assert.Equal(t, addResponse.Device.Revision, response.Device.Revision)
		count++
	}
	assert.Equal(t, 1, count)

	removeResponse, err := client.Remove(context.Background(), &device.RemoveRequest{
		Device: getResponse.Device,
	})
	assert.NoError(t, err)
	assert.NotNil(t, removeResponse)

	eventResponse = <-events
	assert.Equal(t, device.ListResponse_REMOVED, eventResponse.Type)
	assert.Equal(t, device.ID("test1"), eventResponse.Device.ID)
	assert.True(t, eventResponse.Device.Revision > addResponse.Device.Revision)
}

func getDeviceServiceClient() (device.DeviceServiceClient, error) {
	creds, err := getClientCredentials()
	if err != nil {
		return nil, err
	}
	conn, err := grpc.Dial("onos-topo:5150", grpc.WithTransportCredentials(credentials.NewTLS(creds)))
	if err != nil {
		return nil, err
	}
	return device.NewDeviceServiceClient(conn), nil
}

func getClientCredentials() (*tls.Config, error) {
	cert, err := tls.X509KeyPair([]byte(clientCert), []byte(clientKey))
	if err != nil {
		return nil, err
	}
	return &tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
	}, nil
}

const clientCert = `
-----BEGIN CERTIFICATE-----
MIIDZTCCAk0CCQDl7NF6ekffcTANBgkqhkiG9w0BAQsFADByMQswCQYDVQQGEwJV
UzELMAkGA1UECAwCQ0ExEjAQBgNVBAcMCU1lbmxvUGFyazEMMAoGA1UECgwDT05G
MRQwEgYDVQQLDAtFbmdpbmVlcmluZzEeMBwGA1UEAwwVY2Eub3Blbm5ldHdvcmtp
bmcub3JnMB4XDTE5MDQxMTExMTYyM1oXDTIwMDQxMDExMTYyM1owdzELMAkGA1UE
BhMCVVMxCzAJBgNVBAgMAkNBMRIwEAYDVQQHDAlNZW5sb1BhcmsxDDAKBgNVBAoM
A09ORjEUMBIGA1UECwwLRW5naW5lZXJpbmcxIzAhBgNVBAMMGmNsaWVudDEub3Bl
bm5ldHdvcmtpbmcub3JnMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA
5mR12oGXP+uDD7DzQZdTg96eHWTc0UKPwie2I5LLLVsRoH2PO5s2B5r6r/E8OUG4
0pGb6tkDRIJ8eC0Z/6NvBkzn4fsJ5g0UW6sVlXfaf0y9JnMSvV05+g++75a7+CRx
1BG3GNjGWbke1mx8d6SrQ8D1sjI3L0D+32mi0WU9jO2Uw9YXvXgxQmL9Krxdr3M/
aZO9sTJZtIT0EEY3qBpPv+daAbuP5m+uhiEzYZP2bLywyzGyfrUmj9fjG/D1kuMM
haEIUJQ2VTcIApKG/Kb3Mk3b3VCfTvpEHMVrKMoyNHQXXi+6X106+cu2WtoPv+U5
VFVoufjRWSbcOmQ7qIHBiwIDAQABMA0GCSqGSIb3DQEBCwUAA4IBAQBRBR6LTFEU
SWeEeguMsbHxN/6NIZuPejib1q9fTHeZ9cnIHIOLJaZzHiMZn5uw8s6D26kveNps
iCr4O8xOjUa0uwbhMTgm3wkODLlV1DwGjFWk8v5UKGWqUQ94wVMQ16YMIR5DgJJM
0DUzVcoFz+vLnMrDZ0AEk5vra1Z5KweSRvwHX7dJ6FIW7X3IgqXTqJtlV/D/vIi3
UfBnjzqOy2LVfBD7du7i5NbTHfTUpoTvddVwQaKCuQGYHocoQvQD3VQcQDh1u0DD
n2GkeEDLaDAGFAIO+PDg2iT8BhKeEepqswid9gYAhZcOjrlnl6smZo7jEzBj1a9Q
e3q1STjfQqe8
-----END CERTIFICATE-----
`

const clientKey = `
-----BEGIN PRIVATE KEY-----
MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQDmZHXagZc/64MP
sPNBl1OD3p4dZNzRQo/CJ7YjksstWxGgfY87mzYHmvqv8Tw5QbjSkZvq2QNEgnx4
LRn/o28GTOfh+wnmDRRbqxWVd9p/TL0mcxK9XTn6D77vlrv4JHHUEbcY2MZZuR7W
bHx3pKtDwPWyMjcvQP7faaLRZT2M7ZTD1he9eDFCYv0qvF2vcz9pk72xMlm0hPQQ
RjeoGk+/51oBu4/mb66GITNhk/ZsvLDLMbJ+tSaP1+Mb8PWS4wyFoQhQlDZVNwgC
kob8pvcyTdvdUJ9O+kQcxWsoyjI0dBdeL7pfXTr5y7Za2g+/5TlUVWi5+NFZJtw6
ZDuogcGLAgMBAAECggEBAIc9VUjsZSJqVsaxMjnAYI+578qFWHGlxslLkkkTdByt
po005w0wMOkJ+jmpO5bIk3tXadTTim1+wx2wK+C5yQRDxKIMQGVALEEbDlJsxl+P
ZkDZr5hkzxGQiJ4PN0uT6RV5SKdXKCem2Qk5KV751GazMAZoH6inWHVAhwiviw/b
kSJmXcQifxB9R5Br+yCdkRNGg+EtadxAkRtZdW0N0H6LwWxsl32I4o1WM3N2Tyag
kpKPPZ5J5U+279Rpz7W4JAbGzWBOL0Wc2pz5p+aKVTWia0MoqzHR4P4YnkGM+w9Y
j6+Nemdedx62KPhOnQH1uvuG3vnOtt2Ss5OLxePgmjECgYEA9bVguF1D5rpp6MSK
2izZt0mNqhiozm84W2UrAwDhtW5tptW2JBPj2T05+PbEOUEgsvucWfmhZoBXNRCw
IlLQZh46LJFXyW1Awn3PuYquruF61phDoqU9Ou5skJrh0ez+vX872HkH4KW3MfWq
w3LW4qXt6z+lBgPY8hNAlis3WE0CgYEA8Ara5J915ZoVll1As84H61NHmkyMFENh
PjUJqL6tPxvZ+lkBeA157o6mrIgNmG5bLnzonpT4rqemewxEYL39sJ6CVzHRFy8I
F0VNLzZbYizrPLRvT+Gkh0jf6W7Iarzmcdb8cMDxQ+9LmwR/Q3XAD8ntqzrbwVl5
FOZlGq2ZbTcCgYEAuMULlbi07hXyvNLH4+dkVXufZ3EhyBNFGx2J6blJAkmndZUy
YhD+/4cWSE0xJCkAsPebDOI26EDM05/YBAeopZJHhupJTLS2xUsc4VcTo3j2Cdf4
zJ9b2yweQePmuxlwOwop89CYBuw3Rf+KyW1bgJbswkJbE5njE688m3CmLuUCgYAf
K2mtEj++5rky4z0JnBFPL2s20AXIg89WwpBUhx37+ePeLDySmD1jCsb91FTfnETe
zn1uSi3YkBCAHeGrJkCQ9KQ8Kk3aUtMcInWZUdef8fFB2rQxjT1OC9p3d1ky8wCB
e8cf5Q3vIl2Q7Y6Q9fNQmYnxGB19B98/JYOvaSdpFQKBgFBJ+tdJ5ghXSdvAzGno
trQlL1AYW/kYsxZaALd1R+vK3vxeHOtUWiq3923QttYsVXPRQe1TEEdxlOb7+hwE
g5NVOIsDpB1OqjQRb9PjipANkHQRKgrYFB20ZQUoaOMckhlVyqE6WcanGpUxJ0xg
1F0itWrqPGEs83BRQI/aLlsj
-----END PRIVATE KEY-----
`
