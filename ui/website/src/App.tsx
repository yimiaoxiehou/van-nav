import "./App.css";
import Content from "./components/Content";
import DarkSwitch from "./components/DarkSwitch";
import GithubLink from "./components/GithubLink";
import {ToastContainer} from "react-toastify";
function App() {
  return (
    <div className="App">
        <ToastContainer
        position="top-right"
        autoClose={5000}
        hideProgressBar={false}
        newestOnTop={false}
        closeOnClick
        rtl={false}
        pauseOnFocusLoss
        draggable
        pauseOnHover
        theme="light"
    />
        <ToastContainer />
      <DarkSwitch />
      <GithubLink />
      <div className="main">
        <Content />
      </div>
    </div>
  );
}

export default App;
