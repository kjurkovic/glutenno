import React, { useEffect, useState } from 'react'
import { getUserId } from '../../utils/localstorage'

type DeleteFn = (id: string) => void

interface ICommentItem {
  id: string,
  text: string,
  resourceId: string,
  user: string,
  userId: string,
  onDelete: DeleteFn,
}

const CommentItem = (props : ICommentItem) => {

  const [isUserComment, setUserComment] = useState(false)

  const handleDelete = () => props.onDelete(props.id)

  useEffect(() => {
    let userId = getUserId()
    setUserComment(userId === props.userId)
  }, [])

  return (
    <div className="mt-2 pb-2 border-b-2 border-slate-300	">
      <p className="text-base">{props.text}</p>
      <div className="flex flex-row mt-2">
        <h4 className="text-sm font-medium">-- {props.user}</h4>

         { isUserComment ? <button className="ml-6 text-sm text-red-700 hover:text-red-500" onClick={handleDelete}>Delete</button> : <></> }
      </div>
    </div>
  );
}

export default CommentItem;