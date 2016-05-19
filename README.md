Go CAS Client
=============

Forked from : github.com/lucasuyezu/golang-cas-client

WIP to test Jasig CAS server


How to request a Service Ticket on REST API
-------------------------------------------

    import (
      "fmt"
      "github.com/yvesago/golang-cas-client"
    )

    func main() {
      cas := cas.NewClient("https://server", "user", "pass")
      ticket, _ := cas.RequestServiceTicket("service")

      fmt.Println("ST is ", ticket)
    }


How to validate a Service Ticket
-------------------------------

    import (
      "fmt"
      "github.com/yvesago/golang-cas-client"
    )

    func main() {
      cas := cas.NewService("https://server", "service-host")
      response, _ := cas.ValidateServiceTicket("service")

      fmt.Println("ST is ", response.Status)
    }


TODO
----

* Improve error handling
* Invalidate a Service Ticket
* Reuse a TGT (Ticket Granting Ticket) to generate more than one Service Ticket
* Invalidate a TGT (sign out from all services)
