import { commentsClient }  from './client'
import { Comment } from './types';

const commentApi = {
  get: (recipeId: string) => commentsClient.get(`/comments/${recipeId}`, {
    validateStatus: (status) => status === 200
  }).then(res => res.data as Array<Comment>),
  getByResourceOwner: () => commentsClient.get('/comments/user', {
    validateStatus: (status) => status === 200
  }).then(res => res.data as Array<Comment>), 
  delete: (commentId: string) => commentsClient.delete(`/comments/${commentId}`, {
    validateStatus: (status) => status === 204
  }),
  save: (recipeId: string, comment: string, resourceOwnerId: string) => commentsClient.post('/comments', {
    text: comment,
    resourceId: recipeId,
    resourceOwnerId: resourceOwnerId,
  }, {validateStatus: (status) => status === 200 })
    .then(res => res.data)
}

export default commentApi;