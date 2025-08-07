import axios from "axios";

const instance = axios.create({
  baseURL: __API_URL__,
  timeout: 15000,
});

export default instance;
