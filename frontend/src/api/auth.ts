import { setAuthTokens, setUserId } from '../utils/localstorage';
import { authClient }  from './client'
import userApi from './user';

const authApi = {
  login: (email: string, password: string) => authClient.post('/auth/login', {
    email: email,
    password: password,
  }, { validateStatus: (status) => status === 200 }) 
  .then( res => {
    setAuthTokens(res.data.accessToken, res.data.refreshToken)
    return res
  })
  .then(_ => userApi.get()),
  register: (name:string, email: string, password: string) => authClient.post('/auth/register', {
    name: name,
    email: email,
    password: password,
  }, { validateStatus: (status) => status === 200 }) 
  .then( res => {
    setAuthTokens(res.data.accessToken, res.data.refreshToken)
    return res
  })
}

export default authApi;
