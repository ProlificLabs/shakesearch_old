import Loader from '../assets/loading.svg';

const Loading = () => {
  return <div className="flex justify-center my-20 animate-pulse">
    <img src={Loader} className="bg-indigo" alt="" />
  </div>
}

export default Loading;
