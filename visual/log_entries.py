class LogEntry():
    def build_map(self, raw_msg):
        entries = raw_msg.split(" ")
        mp = {}
        for entry in entries:
            key = entry.split("(", 1)[0]
            value = entry.split("(", 1)[1][:-1]
            mp[key] = value
        return mp

    def __init__(self, raw_entry):
        raw_fields = raw_entry.replace('\t', ' ').split(" ", 2)

        self.timestamp = int(raw_fields[0])
        self.type = raw_fields[1]
        self.mp = self.build_map(raw_fields[2])

    def __repr__(self):
        return "<{} TYPE:{}>".format(self.timestamp, self.type)

class JoinEntry(LogEntry):
    def __init__(self, raw_entry):
        super(JoinEntry, self).__init__(raw_entry)
        assert(self.type == "<join>")

        self.nodeId = self.mp["nodeId"]

    def __repr__(self):
        return "<{} JOIN nodeId({})>".format(self.timestamp, self.nodeId)

class QueryEntry(LogEntry):
    def __init__(self, raw_entry):
        super(QueryEntry, self).__init__(raw_entry)
        assert(self.type == "<query>")

        self.key = str(self.mp["key"])
        self.size = int(self.mp["size"])
        self.node = self.mp["node"]
        self.store = self.mp["store"] == "true"

    def __repr__(self):
        return "<{} QUERY key({}) size({}) node({}) store({})>".format(
            self.timestamp, self.key, self.size, self.node, self.store)

class UnderlayPacketEntry(LogEntry):
    def __init__(self, raw_entry):
        super(UnderlayPacketEntry, self).__init__(raw_entry)
        assert(self.type == "<underlay_packet>")

        self.src =  self.mp["src"]
        self.dest = self.mp["dest"]
        if len(self.mp["recv"]) == 0:
            self.recv = None
        else:
            self.recv = self.mp["recv"]
        self.domain = self.mp["domain"]

    def __repr__(self):
        return "<{} UNDERLAY PKT src({}) dest({}) domain({})>".format(
            self.timestamp, self.src, self.dest, self.domain)

class UnderlaySendPacketEntry(UnderlayPacketEntry):
    def __init__(self, raw_entry):
        super(UnderlaySendPacketEntry, self).__init__(raw_entry)
        assert(self.src == self.recv)

    def __repr__(self):
        return "<{} UNDERLAY SEND src({}) dest({})>".format(
            self.timestamp, self.src, self.dest)

class UnderlayRecvPacketEntry(UnderlayPacketEntry):
    def __init__(self, raw_entry):
        super(UnderlayRecvPacketEntry, self).__init__(raw_entry)
        assert(self.dest == self.recv)

    def __repr__(self):
        return "<{} UNDERLAY RECV src({}) dest({})>".format(
            self.timestamp, self.src, self.dest)
