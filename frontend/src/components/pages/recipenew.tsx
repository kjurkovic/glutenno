import { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import recipeApi from "../../api/recipes";
import { RecipeNewType, RecipeType } from "../../api/types";
import { getUserId } from "../../utils/localstorage";

const RecipeNew = () => {
  let navigate = useNavigate()
  
  const [recipe, setRecipe] = useState<RecipeNewType>({
    title: '',
    description: '',
    steps: [],
    ownerId: getUserId() || "",
    views: 0,
  })
  let { recipeId } = useParams();

  const handleChange = (event: any) => {
    let updated = Object.assign({}, recipe);
    switch (event.target.name) {
      case 'title':
        updated.title = event.target.value
        break;
      case 'description':
        updated.description = event.target.value
        break;
    }
    
    setRecipe(updated)
  }

  const handleStep = (event: any) => {
    let updated = Object.assign({}, recipe);
    let step = event.target.name.split('-')[1]
    updated.steps[step - 1].description = event.target.value
    setRecipe(updated)
  }

  const addStepHandler = () => {
    let updated = Object.assign({}, recipe);
    updated.steps.push({description: '', order: recipe.steps.length + 1})
    setRecipe(updated)
  }

  const saveRecipe = () => {
    if (recipeId) {
      recipeApi.update(recipeId, recipe).then(_ => {
        window.alert("Recipe updated")
        navigate(-1)
      })
    } else {
      recipeApi.save(recipe).then(_ => {
        window.alert("Recipe saved")
        navigate(-1)
      })
    }
  }

  useEffect(() => {
    if (recipeId) {
      recipeApi.getById(recipeId).then((item) => setRecipe(item))
    }
  }, [recipeId])
  

  return (
    <div className="container grid grid-cols-2">
      <div>
        <div>
          <label htmlFor="title" className="block text-sm font-medium text-gray-700">
            Title
          </label>
          <div className="mt-1">
            <input
              type="text"
              name="title"
              id="title"
              className="shadow-sm focus:ring-slate-500 focus:border-slate-500 block w-full sm:text-sm border-gray-300 rounded-md p-3"
              placeholder="Enter title"
              value={recipe.title}
              onChange={handleChange}
            />
          </div>
          <div>
            <label htmlFor="description" className="block text-sm font-medium text-gray-700  mt-4">
              Description
            </label>
            <div className="mt-1">
              <textarea
                rows={4}
                name="description"
                id="description"
                className="shadow-sm focus:ring-slate-500 focus:border-slate-500 block w-full sm:text-sm border-gray-300 rounded-md p-4"
                value={recipe.description}
                onChange={handleChange}
              />
            </div>
          </div>
        </div>
         
        <div className="flex flex-row mt-8">
          <h3 className="mt-10 text-2xl font-bold grow">Steps</h3>
          <button type="button" className="bg-amber-600 rounded text-white font-bold hover:bg-amber-400 text-sm px-6 py-0" onClick={addStepHandler}>Add Step</button>
        </div>
        
        <div className="mt-6">
           { 
            recipe.steps
              .sort((first, second) => first.order > second.order ? 1 : -1)
              .map(step => {
                return (
                  <div key={step.order}>
                    <label htmlFor="step" className="block text-sm font-medium text-gray-700  mt-4">
                      Step {step.order}
                    </label>
                    <div className="mt-1">
                      <textarea
                        rows={2}
                        name={`step-${step.order}`}
                        id="step"
                        className="shadow-sm focus:ring-slate-500 focus:border-slate-500 block w-full sm:text-sm border-gray-300 rounded-md p-4"
                        value={step.description}
                        onChange={handleStep}
                      />
                    </div>
                  </div>
                )
              })
          }
          <button type="button" className="bg-amber-600 rounded text-white font-bold hover:bg-amber-400 text-sm p-2 mt-4" onClick={saveRecipe}>Save</button>
        </div>
      </div>  
    </div>
  );
}

export default RecipeNew