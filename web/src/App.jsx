import { lazy } from "solid-js"
import { Router, Route } from "@solidjs/router"
import { Toaster } from 'solid-toast'

import NavBar from "./NavBar"
import Main from "./Main"
const Login = lazy(() => import("./Login"))
const Register = lazy(() => import("./Register"))
const NotFound = lazy(() => import("./NotFound"))
const Logout = lazy(() => import("./modal/Logout"))
function App() {
  return (
    <div>
      <NavBar />
      <Router>
        <Route path="/" component={Main}/>
        <Route path="/login" component={Login} />
        <Route path="/logout" component={Logout} />
        <Route path="/register" component={Register} />
        {/* <Route path="/registry">
          <Route path="/" component={Registry} />
          <Route path="/:name" component={Repo} />
          <Route path="/:name/:image" component={Image} />
        </Route> */}
        {/* <Route path="/settings" component={Settings} /> */}
        <Route path="*" component={NotFound} />
      </Router>
      <Toaster />
    </div>
  )
}

export default App
