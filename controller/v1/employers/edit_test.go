package employers_test

// func Test_Employer_EditJob_IncorrectMethod(t *testing.T) {
// 	assert := assert.New(t)
// 	ts := httptest.NewServer(http.HandlerFunc(employers.EditJob))

// 	defer ts.Close()

// 	methods := []string{"GET", "PUT", "DELETE"}

// 	for _, method := range methods {

// 		request, err := http.NewRequest(method, ts.URL, nil)

// 		if err != nil {
// 			t.Fatal()
// 		}

// 		httpClient := &http.Client{}

// 		response, err := httpClient.Do(request)

// 		assert.Nil(err)
// 		assert.Equal(int(http.StatusMethodNotAllowed), response.StatusCode)

// 	}

// }

// func Test_Employer_EditJob_Correct_NoData(t *testing.T) {
// 	assert := assert.New(t)
// 	ts := httptest.NewServer(http.HandlerFunc(employers.EditJob))

// 	defer ts.Close()

// 	employer := testhelper.Helper_RandomEmployer(t)
// 	job := testhelper.Helper_RandomJob(employer, t)

// 	data := map[string]string{
// 		"publicid":      job.PublicID,
// 		"title":         "",
// 		"streetaddress": "",
// 		"city":          "",
// 		"zipcode":       "",
// 		"tags":          "",
// 		"description":   "",
// 	}

// 	requestBody, err := json.Marshal(data)

// 	if err != nil {
// 		t.Fatal()
// 	}

// 	request, err := http.NewRequest("POST", ts.URL, bytes.NewBuffer(requestBody))

// 	if err != nil {
// 		t.Fatal()
// 	}

// 	token, err := jwt.GenerateToken(employer.PublicID)

// 	if err != nil {
// 		t.Fatal()
// 	}

// 	request.Header.Set("Content-Type", "application/json")
// 	request.Header.Set("Authorization", "Bearer "+token)

// 	httpClient := &http.Client{}

// 	response, err := httpClient.Do(request)

// 	assert.Nil(err)
// 	assert.NotNil(response)

// }

// func Test_Employer_EditJob_Correct_Data(t *testing.T) {
// 	assert := assert.New(t)
// 	ts := httptest.NewServer(http.HandlerFunc(employers.EditJob))

// 	defer ts.Close()

// 	employer := testhelper.Helper_RandomEmployer(t)
// 	job := testhelper.Helper_RandomJob(employer, t)

// 	data := map[string]string{
// 		"publicid":      job.PublicID,
// 		"title":         "A New Job",
// 		"streetaddress": "123 Address",
// 		"city":          "City",
// 		"zipcode":       "00000",
// 		"tags":          "full-time,remote",
// 		"description":   "this is a job",
// 	}

// 	requestBody, err := json.Marshal(data)

// 	if err != nil {
// 		t.Fatal()
// 	}

// 	request, err := http.NewRequest("POST", ts.URL, bytes.NewBuffer(requestBody))

// 	if err != nil {
// 		t.Fatal()
// 	}

// 	token, err := jwt.GenerateToken(employer.PublicID)

// 	if err != nil {
// 		t.Fatal()
// 	}

// 	request.Header.Set("Content-Type", "application/json")
// 	request.Header.Set("Authorization", "Bearer "+token)

// 	httpClient := &http.Client{}

// 	response, err := httpClient.Do(request)

// 	var result map[string]string
// 	decoder := json.NewDecoder(response.Body)
// 	decoder.Decode(&result)

// 	assert.Nil(err)
// 	assert.Equal(data["title"], result["title"])

// }
