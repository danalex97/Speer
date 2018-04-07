package examples
//
// import (
//   "github.com/danalex97/Speer/interfaces"
//   "runtime"
//   "sync"
//   "fmt"
// )
//
// type SimpleTorrent struct {
//   sync.Mutex
//
//   id    string
//   ids   []string
//   links map[string]interfaces.Link
// }
//
// type idBroadcast struct {
//   ids []string
// }
//
// /* Interface functions. */
// func (s *SimpleTorrent) OnJoin() {
//   go func() {
//     for {
//       if s.Transfer() == nil {
//         // engine not ready
//
//         runtime.Gosched()
//         continue
//       }
//
//       // check links
//       for _, l := range s.links {
//         select {
//         case m, _ := <-l.Download():
//             // data := m.(Data)
//             fmt.Println(s.id, "data", m)
//         default:
//           continue
//         }
//       }
//
//       // check for control messages
//       select {
//       case _, ok := <-s.UnreliableNode().Recv():
//         if ok {
//           fmt.Println("Packet receive!")
//         }
//
//       case m, ok := <-s.Transfer().ControlRecv():
//         if !ok {
//           continue
//         }
//
//         switch msg := m.(type) {
//         case idBroadcast:
//           s.updateIds(msg.ids)
//           fmt.Println(s.id, "received", msg.ids)
//         }
//
//       default:
//         runtime.Gosched()
//       }
//     }
//   }()
//
//   go func() {
//     // wait for engine to be ready
//     Wait(func () bool {
//       return s.Transfer() == nil
//     })
//
//     // broadcast neighbours
//     for _, id := range s.ids {
//       if id != s.id {
//         Wait(func () bool {
//           return !s.Transfer().ControlPing(id)
//         })
//
//         s.Transfer().ControlSend(id, idBroadcast{s.ids})
//       }
//     }
//   }()
// }
//
// func (s *SimpleTorrent) OnQuery(query interfaces.Query) error {
//   return nil
// }
//
// func (s *SimpleTorrent) OnLeave() {
// }
//
// func (s *SimpleTorrent) New() DHTNode {
//   // Constructor that assumes the UnreliableNode component is filled in
//   node := new(SimpleTorrent)
//
//   node.id    = node.UnreliableNode().Id()
//   node.ids   = []string{node.id, node.UnreliableNode().Join()}
//   node.links = map[string]interfaces.Link{}
//
//   return node
// }
//
// func (s *SimpleTorrent) Key() string {
//   return RandomKey()
// }
//
// /* Local functions */
// func (s *SimpleTorrent) updateIds(ids []string) {
//   allIds := make(map[string]bool)
//   for _, id := range ids {
//     allIds[id] = true
//   }
//   for _, id := range s.ids {
//     allIds[id] = true
//   }
//
//   s.ids = []string{}
//   for id, _ := range allIds {
//     s.ids = append(s.ids, id)
//
//     if id == s.id {
//       continue
//     }
//
//     // register link if not registered
//     if _, ok := s.links[id]; !ok {
//       s.links[id] = s.Transfer().Connect(id)
//
//       // if the link is new, we broadcast our list again
//       Wait(func () bool {
//         return !s.Transfer().ControlPing(id)
//       })
//       s.Transfer().ControlSend(id, idBroadcast{s.ids})
//
//       // send a big packet
//       s.links[id].Upload(interfaces.Data{s.Key(), 1000})
//     }
//   }
// }
