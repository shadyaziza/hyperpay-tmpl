package main

import (
	"fmt"
	"html/template"
	"net/http"
)

type templateData struct {
	CheckoutID        string
	CheckoutURL       string
	CustomCheckoutURL string
	RedirectURL       string
	Total             float64
}

func main() {
	filerServer := http.FileServer(http.Dir("./ui/static/"))
	http.Handle("/static/", http.StripPrefix("/static", filerServer))
	http.HandleFunc("/arco/payment", paymentHandler)
	http.ListenAndServe(":4141", nil)
}

func paymentHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("checkoutId")

	if id == "" {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized operation, please make sure you have not sent an empty checkout id"))
		return
	}

	files := []string{"./ui/payment.tmpl.html", "./ui/redirect.tmpl.html"}

	t, err := template.ParseFiles(files...)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Unable to load payment page, please try again later"))
		return
	}

	data := &templateData{
		CheckoutID:        id,
		CheckoutURL:       fmt.Sprintf(`https://eu-test.oppwa.com/v1/paymentWidgets.js?checkoutId=%s`, id),
		CustomCheckoutURL: fmt.Sprintf(`https://eu-test.oppwa.com/v1/checkouts/%s/payment`, id),
		RedirectURL:       `https://arcopayment.web.app`,

		Total: 1402.12,
	}
	err = t.ExecuteTemplate(w, "payment", data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Unable to load payment page, please try again later"))
		return
	}

}
