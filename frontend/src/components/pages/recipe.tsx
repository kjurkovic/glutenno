import React, { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import recipeApi from "../../api/recipes";
import commentApi from "../../api/comments";
import { RecipeType, Comment } from "../../api/types";
import { getAccessToken, getUserId } from "../../utils/localstorage";
import Step from "../shared/step";
import CommentItem from "../shared/commentitem";

const Recipe = () => {
  const [recipe, setRecipe] = useState<RecipeType>()
  const [comments, setComments] = useState<Array<Comment>>([])
  let { recipeId } = useParams();

  const [canComment, setCanComment] = useState(false)
  const [comment, setComment] = useState("")

  const loadComments = () => {
    commentApi.get(recipeId as string).then((data) => setComments(data))
  }

  useEffect(() => {
    recipeApi.view(recipeId || '')
    setCanComment(getAccessToken() != null)
    recipeApi.getById(recipeId as string).then((data) => setRecipe(data))
    loadComments()
  }, [])

  const handleCommentChange = (event: any) => {
    setComment(event.target.value)
  }

  const submitComment = () => {
    let userId = getUserId()

    if (comment.trim().length == 0 || recipe === undefined || userId === null) return

    commentApi.save(recipe.id, comment, userId).then(() => {
      setComment('')
      loadComments()
    })
  }

  const handleCommentDelete = (id: string) => {
    if (window.confirm(`Are you sure you want to delete this comment?`)) {
       commentApi.delete(id).then(_ => {
        let filtered = comments.filter(item => item.id !== id)
        setComments(filtered)
       })
    } 
  }

  if (recipe) {
    return (
      <div className="container">
        <h1 className="text-5xl font-bold">{recipe.title}</h1>
        <p className="mt-6">{recipe.description}</p>

        <h3 className="mt-10 text-2xl font-bold">Steps</h3>
        <div className="mt-6">
          { 
            recipe.steps
              .sort((first, second) => first.order > second.order ? 1 : -1)
              .map(step => <Step key={step.id} {...step} />)
          }
        </div>
        <div>
          { canComment ? 
            <div className="my-10">
              <label htmlFor="comment" className="block text-sm font-medium text-gray-700">
                Add your comment
              </label>
              <div className="mt-1">
                <textarea
                  rows={4}
                  name="comment"
                  id="comment"
                  className="shadow-sm p-2 focus:ring-slate-500 focus:border-slate-500 block w-full sm:text-sm border-gray-300 rounded-md"
                  value={comment}
                  onChange={handleCommentChange}
                />
                <button type="button" className="-inline-flex mt-2 items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-slate-600 hover:bg-slate-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-slate-500" onClick={submitComment}>Post</button>
              </div>
            </div>
            : <></> 
          }
          <div>
            <h2 className="text-lg font-bold">Comments</h2>
            {
              comments.map((item) => <CommentItem key={item.id} {...item} onDelete={handleCommentDelete} />)
            }
          </div>
        </div>
      </div>
    );
  }
  return <></>
}

export default Recipe;
