
function ServeDns(ctx, w, r)
        for i, q in r.Question() do
                printf("q %v\n", q)
                ip = net.ParseIP("192.168.6.2")
                name = q.Name
                reply = ctx:ReplyA(r, name, ip, default_ttl)
                printf("reply %v\n", reply)
                w:WriteMsg(reply)
                return
        end
end

