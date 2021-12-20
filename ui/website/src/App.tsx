import './App.scss';
import Header from './components/Header';
import Content from './components/Content';
import About from "./components/About";
// import Footer from './components/Footer';
// import data from './data';
import FetchList from './utils/api';
import { useState, useEffect,useCallback } from 'react';

function App() {
  const [curr, setCurr] = useState<string>();
  const [currData,setCurrData] = useState([]);
  const dealData = (old) => {
    // 把老的无序数据洗一下
    const newData = [];
    const type1 = old.map((item) => {
      return item.catelog;
    });
    const newTypes = Array.from(new Set(type1));
    // 遍历新的类型，然后挨个找
    for (let item of newTypes) {
      const newType = {type: item};
      const l = [];
      // 搜索全部的数据找
      for (let olditem of old){
        if (olditem.catelog === item){
          l.push(olditem);
        }
      }
      newType['list'] = l;
      newData.push(newType);
    }
    return newData;
  };
  const loadData = useCallback(async()=>{
    try {
      const r = await FetchList();
      const newData = dealData(r?.data?.tools);
      if (!curr && newData.length) {
        setCurr(newData[0]?.type)
      }
      // 判断数据
      if(window.localStorage.getItem('door') ){
        // 改成根据字段进行判断了
        setCurrData(newData);
      } else {
        setCurrData(dealData(r?.data?.tools?.filter((item) => {
          if (r.data && r.data.privateList && r.data.privateList.includes(item.catelog)) {
            return false;
          } else {
            return true;
          }
        })))
      }
    } catch (e) {
      console.log(e);
    }
  },[setCurrData,curr,setCurr]);
  useEffect(()=>{
    if (window.location.search === "?door"){
      window.localStorage.setItem("door","1")}
    if (currData.length === 0) {
      loadData();
    }
    if (window.localStorage.getItem("curr")){
      setCurr(window.localStorage.getItem("curr"));
    }


  },[currData,loadData]);


  const handleChangeHeader = (newType: string) => {
    window.localStorage.setItem("curr",newType);
    setCurr(newType);
  }

  const renderMain = () => {
    if (currData.length){
      return (
        <Content data={currData} curr={curr} />
      );
    }
  };
  return (
    <div className="App" >
      <Header onClick={handleChangeHeader} data={currData} curr={curr} />
      <div className="main">
        {renderMain()}
      </div>
    </div>
  );
}

export default App;
