import { setAuthTokens, setUserId } from '../utils/localstorage';
import { userClient }  from './client'

const userApi = {
  get: () => userClient.get('/user', { validateStatus: (status) => status === 200 }) 
  .then( res => {
    setUserId(res.data.id)
    return res
  }),
}

export default userApi;
