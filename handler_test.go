package main

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestEdiToJsonHandler(t *testing.T) {
	// Create a sample .edi file for the test
	testCases := []struct {
		name     string
		fileName string
		content  string
		wantCode int
	}{
		{
			name:     "277",
			fileName: "sample.edi",
			content: `ISA*00*          *00*          *ZZ*RECEIVERID     *ZZ*SENDERID       *210101*1205*U*00401*000000002*0*P*>
GS*HN*RECEIVERID*SENDERID*20210101*1205*2*X*004010X093A1
ST*277*0002*004010X093A1
BHT*0010*11*20004567*20210101*1205
HL*1**20*1
NM1*PR*2*PAYERNAME*****PI*12345
HL*2*1*21*1
NM1*41*2*RECEIVER COMPANY*****46*987654321
HL*3*2*19*1
NM1*1P*1*DOCTOR NAME*****XX*1234567893
HL*4*3*PT*0
NM1*QC*1*PATIENT NAME****MI*123456789A
TRN*2*XYZ12345ABC*9876543210
STC*A1:20:PR*20210101*WQ*1234.56
REF*D9*55555
DTP*472*D8*20201215
SE*16*0002
GE*1*2
IEA*1*000000002`,
		},
		{
			name:     "276",
			fileName: "sample.edi",
			content: `ISA*00*          *00*          *ZZ*SENDERID       *ZZ*RECEIVERID     *210101*1200*U*00401*000000001*0*P*>
GS*HS*SENDERID*RECEIVERID*20210101*1200*1*X*004010X093A1
ST*276*0001*004010X093A1
BHT*0010*13*10001234*20210101*1200
HL*1**20*1
NM1*PR*2*PAYERNAME*****PI*12345
HL*2*1*21*1
NM1*41*2*SENDER COMPANY*****46*987654321
HL*3*2*19*1
NM1*1P*1*DOCTOR NAME*****XX*1234567893
HL*4*3*22*0
TRN*1*ABC12345XYZ*9876543210
REF*D9*55555
DTP*472*D8*20201215
SE*12*0001
GE*1*1
IEA*1*000000001`,
		},
		{
			name:     "216",
			fileName: "sample.edi",
			content: `ISA*00*          *00*          *02*CARRIERID      *ZZ*RECEIVERID     *220101*1500*U*00401*000000001*0*P*>
GS*OG*CARRIERID*RECEIVERID*20220101*1500*1*X*004010
ST*216*0001
B10*1234567890*SHIPMENT123*CAR001
N1*SF*SHIPPER NAME*9*123456789
N3*123 SHIPPER STREET
N4*SHIPPER CITY*ST*12345*US
N1*ST*CONSIGNEE NAME*9*987654321
N3*456 CONSIGNEE AVE
N4*CONSIGNEE CITY*ST*54321*US
G62*37*20220101*11*1500
AT7*X1*NS***20220101*1400
MS1*ORIGIN CITY*ST*US
MS2*CAR001*12345
L11*REF123*BM
L11*PO123456*PO
L11*INV123456*IV
SE*15*0001
GE*1*1
IEA*1*000000001`,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			err := os.WriteFile(tt.fileName, []byte(tt.content), 0644)
			if err != nil {
				t.Fatalf("Failed to create sample .edi file: %v", err)
			}
			defer os.Remove(tt.fileName)

			// Create a buffer to hold the multipart form data
			var buf bytes.Buffer
			writer := multipart.NewWriter(&buf)

			// Add the .edi file to the form data
			fileWriter, err := writer.CreateFormFile("ediFile", tt.fileName)
			if err != nil {
				t.Fatalf("Failed to create form file: %v", err)
			}

			file, err := os.Open(tt.fileName)
			if err != nil {
				t.Fatalf("Failed to open sample .edi file: %v", err)
			}
			defer file.Close()
			_, err = io.Copy(fileWriter, file)
			if err != nil {
				t.Fatalf("Failed to copy file content: %v", err)
			}
			writer.Close()

			// Create a new HTTP request with the form data
			req := httptest.NewRequest("POST", "/edi-to-json", &buf)
			req.Header.Set("Content-Type", writer.FormDataContentType())

			// Create a ResponseRecorder to capture the response
			w := httptest.NewRecorder()

			// Call the handler function
			handler := http.HandlerFunc(ediToJsonHandler)
			handler.ServeHTTP(w, req)

			// Check the response status code
			if status := w.Code; status != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
			}
		})
	}

}
