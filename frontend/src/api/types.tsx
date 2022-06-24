export interface Step {
  id: string,
  description: string,
  order: number,
}

export interface RecipeType {
  id: string,
  title: string,
  description: string,
  steps: Array<Step>,
  ownerId: string,
  views: number,
}

export interface RecipeNewType {
  title: string,
  description: string,
  steps: Array<StepNew>,
  ownerId: string,
  views: number,
}

export interface StepNew {
  description: string,
  order: number,
}


export interface Comment {
  id: string,
  text: string,
  resourceId: string,
  user: string,
  userId: string,
}