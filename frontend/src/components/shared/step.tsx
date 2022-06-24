import React from "react";

interface IStep {
  id: string,
  description: string,
  order: number,
}

const Step = (props: IStep) => {
  return (
    <div className="flex flex-row items-center mb-4">
      <div className="flex-none	p-2 rounded-full w-10 h-10 text-center align-middle bg-amber-400 text-white font-bold">{props.order}</div> 
      <p className="ml-4">{props.description}</p>
    </div>
  )
}

export default Step;
