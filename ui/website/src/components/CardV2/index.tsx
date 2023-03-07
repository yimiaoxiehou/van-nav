import {useMemo} from "react";
import "./index.css";
import copy from "copy-to-clipboard";
import {toast} from 'react-toastify'
import 'react-toastify/dist/ReactToastify.css';


function handlerCopy(value, msg) {
    copy(value);
    toast(msg);
}

const Card = ({title, url, des, logo, catelog, username, password, onClick}) => {
    const el = useMemo(() => {
        if (url === "admin") {
            return <img src={logo} alt={title}/>
        } else {
            if (logo.split(".").pop().includes("svg")) {
                return <embed src={`/api/img?url=${logo}`} type="image/svg+xml"/>
            } else {
                return <img src={`/api/img?url=${logo}`} alt={title}/>
            }
        }
    }, [logo, title, url])
    return (
        <div className="card-box"
        >
            <a href={url}
               onClick={() => {
                   onClick();
               }}
               target="_blank"
               rel="noreferrer"
               style={{height: '80px'}}
               className="card-content">
                <div className="card-left">
                    {el}
                    {/* {url === "admin" ? (
            <img src={logo} alt={title} />
          ) : (
            <img src={`/api/img?url=${logo}`} alt={title} />
          )} */}
                </div>
                <div className="card-right">
                    <div className="card-right-top">
                        <span className="card-right-title" title={title}>{title}</span>
                        <span className="card-tag" title={catelog}>{catelog}</span>
                    </div>
                    <div title={des} className="card-right-bottom">{des}</div>
                </div>
            </a>
            <div className="card-content" style={{height: '35px'}}>
                    <p className="card-right-bottom" style={{width: '45%'}} onClick={() => username !== undefined && username!=="" && handlerCopy(username, "用户复制成功")}>
                        用户：{username}
                    </p>
                    <p className="card-right-bottom" style={{width: '53%', marginLeft: '2%'}} onClick={() => password !== undefined && password !== "" && handlerCopy(password, "密码复制成功")}>
                        密码：{password}
                    </p>
            </div>
        </div>
    );
};

export default Card;
