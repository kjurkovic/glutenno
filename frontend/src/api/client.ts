import axios, { AxiosInstance, AxiosRequestConfig } from "axios";
import { getAccessToken, getRefreshToken, removeTokens, setAuthTokens } from "../utils/localstorage";
import createAuthRefreshInterceptor from 'axios-auth-refresh';

const authClient: AxiosInstance = axios.create({
  baseURL: 'http://localhost:9090',
  timeout: 1000,
  headers: {
    'Access-Control-Allow-Origin': '*',
    'Content-Type': 'application/json'
  }
});

const userClient: AxiosInstance = axios.create({
  baseURL: 'http://localhost:9090',
  timeout: 1000,
  headers: {
    'Access-Control-Allow-Origin': '*',
    'Content-Type': 'application/json'
  }
});

const recipeClient: AxiosInstance = axios.create({
  baseURL: 'http://localhost:9091',
  timeout: 1000,
  headers: {
    'Access-Control-Allow-Origin': '*',
    'Content-Type': 'application/json'
  }
});

const commentsClient: AxiosInstance = axios.create({
  baseURL: 'http://localhost:9092',
  timeout: 1000,
  headers: {
    'Access-Control-Allow-Origin': '*',
    'Content-Type': 'application/json'
  }
});

const accessTokenInterceptor = (request: AxiosRequestConfig<any>) => {
  const token = getAccessToken()
  if (request.headers != undefined && token != null) {
    request.headers['Authorization'] = `Bearer ${token}`;
  }
  return request;
}

const refreshTokenRetrier = (failedRequest: any) => authClient.post('http://localhost:9090/auth/refresh-token', { 
  refreshToken: getRefreshToken() 
}, { validateStatus: (status) => status === 200 }
).then(tokenRefreshResponse => {
  setAuthTokens(tokenRefreshResponse.data.accessToken, tokenRefreshResponse.data.refreshToken)
  failedRequest.response.config.headers['Authorization'] = `Bearer ${tokenRefreshResponse.data.accessToken}`;
  return Promise.resolve();
}).catch((err) => {
  console.log(err);
  removeTokens()
  window.location.assign('/signin');
});

userClient.interceptors.request.use(accessTokenInterceptor)
recipeClient.interceptors.request.use(accessTokenInterceptor)
commentsClient.interceptors.request.use(accessTokenInterceptor)

createAuthRefreshInterceptor(userClient, refreshTokenRetrier);
createAuthRefreshInterceptor(recipeClient, refreshTokenRetrier);
createAuthRefreshInterceptor(commentsClient, refreshTokenRetrier);

export { authClient, userClient, recipeClient, commentsClient };
