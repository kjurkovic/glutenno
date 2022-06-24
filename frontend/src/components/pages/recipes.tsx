import React, { useEffect, useState } from "react";
import recipeApi from '../../api/recipes'
import { RecipeType } from "../../api/types";
import RecipeItem from "../shared/recipeitem";

const Recipes = () => {

  const [recipes, setRecipes] = useState<Array<RecipeType>>([])

  useEffect(() => {
    recipeApi.get().then((data) => setRecipes(data))
  }, [])

  if (recipes.length) {
    return (
      <div className="grid grid-cols-3 gap-2">
        { recipes.map(recipe => <RecipeItem key={recipe.id} {...recipe} />)}
      </div>
    );
  } else {
    return (
      <p>There aren't any recipes available at the moment.</p>
    );
  }
}

export default Recipes;