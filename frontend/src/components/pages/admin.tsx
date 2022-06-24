import React, { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import commentApi from "../../api/comments";
import recipeApi from "../../api/recipes";
import { RecipeType } from "../../api/types";
import { Chart as ChartJS, ArcElement, Tooltip, Legend, ChartData, CategoryScale, LinearScale, BarElement, Title,  } from 'chart.js';
import { Doughnut, Bar } from 'react-chartjs-2';
import randomcolor from "../../utils/randomcolor";
import { Comment } from "../../api/types"

ChartJS.register(ArcElement, Tooltip, Legend,  CategoryScale, LinearScale, BarElement, Title);

type DougnutDataType = Map<string, Comment[]>
type DeleteFn = (id: string) => void

const RecipeRow = (props:{id: string, title: string, delete: DeleteFn}) => {

  const deleteHandler = () => {
    if (window.confirm('Are you sure you want to delete this recipe?')) {
      props.delete(props.id)
    }
  }
  return (
    <div className="my-2 p-3 bg-slate-300 rounded flex flex-row">
    <Link to={`/admin/recipe/${props.id}`} className="text-sm hover:underline grow">
      <div>
        <h2 className="font-medium">{props.title}</h2>
        <p>{props.id}</p>
      </div>
    </Link>
    <button className="text-sm text-red-700" onClick={deleteHandler}>Delete</button>
    </div>
  );
}

const Admin = () => {

  const [recipes, setRecipes] = useState<Array<RecipeType>>([])
  const [commentsData, setCommentsData] = useState({
    labels: [''],
    datasets: [{ label: "", data: [0], backgroundColor: [''] }],
    borderWidth: 1,
  })

  const [recipeViewData, setRecipeViewData] = useState({
    labels: [''],
    datasets: [{ label: "", data: [0], backgroundColor: '' }],
  })

  const barOptions = {
    responsive: true,
    plugins: {
      legend: {
        position: 'top' as const,
      },
      title: {
        display: true,
        text: 'Recipe views',
      },
    },
  };

  const dougnutOptions = {
    responsive: true,
    plugins: {
      legend: {
        position: 'top' as const,
      },
      title: {
        display: true,
        text: 'Comments per recipe',
      },
    },
  };

  const loadRecipes = () => {
    recipeApi.getByUser()
      .then((items) => {
        setRecipes(items)
        return items
      })
      .then(items => {
        setRecipeViewData({
          labels: items.map(item => item.title),
          datasets: [
            {
              label: "Views",
              data: items.map(item => item.views),
              backgroundColor: randomcolor(),
            }
          ],
        })
      })
  }

  const loadCommentStats = () => {
    commentApi.getByResourceOwner()
      .then(res => {
        return res.reduce<DougnutDataType>((acc, item) => {
          acc.set(item.resourceId, acc.get(item.resourceId) ?? [])
          acc.get(item.resourceId)?.push(item)
          return acc
        }, new Map())
      })
      .then(data => {
        let arr = Array.from(data)
        let labels = arr.map(item => item[0])
        setCommentsData({
          labels: labels,
          datasets: [
            {
              label: "Comments per recipe",
              data: arr.map(item => item[1].length),
              backgroundColor: arr.map(_ => randomcolor())
            }
          ],
          borderWidth: 1,
        })
      })
  }

  const deleteHandler = (id: string) => {
    recipeApi.delete(id).then(() => {
      loadRecipes()
      loadCommentStats()
    })
  }

  useEffect(() => {
    loadRecipes()
    loadCommentStats()
  }, [])

  return (
    <div>
      <div className="flex flex-row">
        <h2 className="text-lg font-bold grow">Your recipes</h2>
        <Link to="/admin/recipe/add" className="bg-amber-600 rounded text-white font-bold hover:bg-amber-400 text-sm p-2">New recipe</Link>
      </div>
      <div className="grid grid-cols-3 gap-4">
        <div>
          { recipes.map(item => <RecipeRow key={item.id} {...item} delete={deleteHandler} /> )}
        </div>
        <div>
          <Doughnut options={dougnutOptions} data={commentsData} />
        </div>
        <div>
          <Bar options={barOptions} data={recipeViewData} />
        </div>
      </div>
      
    </div>
  )
}

export default Admin;