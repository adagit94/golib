# Overview

WebSockets for everyday use.

## Examples
Example of bootstrapped connection upgrade logic in context of fasthttp:

    upgrader := ws.FastHTTPUpgrader{
      ReadBufferSize:  1024,
      WriteBufferSize: 1024,
      CheckOrigin:     func(r *fasthttp.RequestCtx) bool { return true },
    }

    hub := wio.CreateBroadHub(wio.CommonClientOptions[uint32]{
      PingInterval:  time.Second * 5,
      PongWait:      time.Second * 2,
      WriteBuffSize: 100,
      MaxReadBytes:  1024,
    })

    var counter atomic.Uint32

    fasthttp.ListenAndServe("localhost:4000", func(ctx *fasthttp.RequestCtx) {
      err := upgrader.Upgrade(ctx, func(conn *ws.Conn) {
        // Upgrade successfull.
        clientId := counter.Load()
        terminated := hub.Subscribe(clientId, conn)

        counter.Add(1)
        <-terminated // Necessary in context of fasthttp to prevent dangling pointer dereference of lost ws.Conn. It doesn't block current execution of the outer fasthttp handler.
        fmt.Printf("Terminated client %v.", clientId)
      })

      if err != nil {
        fmt.Println("Failed to upgrade connection.")
      }
    })