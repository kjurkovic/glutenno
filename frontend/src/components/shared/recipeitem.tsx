import React from "react";
import { Link } from "react-router-dom";

interface IStep {
  id: string,
  description: string,
  order: number,
}

interface IRecipeItem {
  id: string,
  title: string,
  description: string,
  steps: Array<IStep>
}

const RecipeItem = (props: IRecipeItem) => {
  return (
    <div className="bg-stone-100 px-4 py-5 mr-4 mb-4 shadow-sm rounded-lg min-w-[40%]">
      <div className="text-[#333333]">
          <h1 className="text-xl font-bold">{props.title}</h1>
          <div className="mt-6">{props.description}</div>
          <div className="mt-6 text-right">
            <Link to={`/recipes/${props.id}`} className="mt-6 hover:underline">Cook it</Link>
          </div>
      </div>
    </div>
  );
}

export default RecipeItem;