import { recipeClient }  from './client'
import { RecipeNewType, RecipeType } from './types';

const recipeApi = {
  get: () => recipeClient.get('/recipes', {
    validateStatus: (status) => status === 200
  }).then(res => res.data as Array<RecipeType>),
  getByUser: () => recipeClient.get('/recipes/user', {
    validateStatus: (status) => status === 200
  }).then(res => res.data as Array<RecipeType>),
  getById: (id: string) => recipeClient.get(`/recipes/${id}`, {
    validateStatus: (status) => status === 200
  }).then(res => res.data as RecipeType),
  view: (id: string) => recipeClient.put(`/recipes/${id}/view`, {}, {
    validateStatus: (status) => status === 200
  }),
  delete: (id: string) => recipeClient.delete(`/recipes/${id}`, {
    validateStatus: (status) => status === 204
  }), 
  update: (id: string, recipe: RecipeNewType) => recipeClient.put(`/recipes/${id}`, recipe, {
    validateStatus: (status) => status === 200
  }),
  save: (recipe: RecipeNewType) => recipeClient.post("/recipes", recipe, {
    validateStatus: (status) => status === 200
  }),
}

export default recipeApi;