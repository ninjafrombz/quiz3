module Todo exposing (main)

import Browser
import Html exposing (..)
import Html.Attributes exposing (..)
import Html.Events exposing (onClick, onInput, onSubmit)


type alias Model =
    { field : String
    , uid : Int
    , todos : List Todo
    }


type alias Todo =
    { id : Int
     , task_name : String
     , description : String
     , category : String
     , priority : String
     , isComplete: Bool
    }

type Msg
    = AddTodo
    | SetField String
    | DeleteTodo Int
    | CompleteTodo Int Bool


initialModel : Model
initialModel =
    { field = ""
    , uid = 0
    , todos = []
    }

view : Model -> Html Msg
view model =
    div [ class "text-center" ] [ div [ class "todo-container text-left p-24 bg-white shadow-sm rounded flex flex-col mx-auto my-48" ]
        [ header [ ]
            [ h1 [ class "text-24 font-bold mb-24" ] [ text " Desire's Todo List" ]
            ]
        , Html.form [ class "w-full flex justify-between" ,onSubmit AddTodo ] [
            input
                [ class "todo-input"
                , onInput
                    (\string -> SetField string)
                , value model.field
                ]
                []
            , button [ class "btn", type_ "submit", disabled (model.field == "") ] [ text "Create" ]
        ]
        , ul [ class "text-left mt-24" ] (List.map viewSearchResult model.todos)
    ]
    , a [ class "mt-48 text-12 text-gray-800", href "https://github.com/ninjafrombz/quiz3" ] [ text "Github" ]
    ]


viewSearchResult : Todo -> Html Msg
viewSearchResult todo =
    li [ class "border-b border-gray-200 py-8 flex justify-between", onClick (CompleteTodo todo.id todo.isComplete) ]
        [ span [ classList[("completed", todo.isComplete)], class "text-todo" ] [ text todo.task_name ]
        , button
            [ class "text-gray-800 outline-none", onClick (DeleteTodo todo.id)]
            [ text "X" ]
        ]

update : Msg -> Model -> Model
update msg model =
    case msg of
        AddTodo ->
            { model | todos = { id = model.uid, task_name = model.field, description = model.field, category = model.field, priority = model.field, isComplete = False } :: model.todos, field = "", uid = model.uid + 1 }
        SetField str ->
            { model | field = str }
        DeleteTodo id ->
            { model | todos = List.filter(\todo -> todo.id /= id) model.todos }
        CompleteTodo id complete ->
            let
                updateTodo todo =
                    if todo.id == id then
                        { todo | isComplete = not complete }
                    else
                        todo
            in
            { model | todos = List.map updateTodo model.todos }



main : Program () Model Msg
main =
    Browser.sandbox
        { view = view
        , update = update
        , init = initialModel
        }







































-- module Todo exposing (main)
-- import Browser
-- import Html exposing (..)
-- import Html.Attributes exposing (..)
-- import Html.Events exposing (onInput)



-- -- MAIN
-- main =
--     Browser.sandbox { init = init, update = update, view = view }

-- -- MODEL

-- type alias Model = 
--     { task_name : String
--   , description : String
--   , category : String
--   , priority : String
--   , status : String
--     }

-- init : Model
-- init =
--   Model  "" "" "" "" ""



-- -- UPDATE


-- type Msg
--   = Task_Name String
--   | Description String
--   | Category String
--   | Priority String
--   | Status String


-- update : Msg -> Model -> Model
-- update msg model =
--   case msg of
--     Task_Name task_name ->
--       { model | task_name = task_name }

--     Description description ->
--       { model | description = description }

--     Category category ->
--       { model | category = category }

--     Priority priority ->
--       { model | priority = priority }

--     Status status ->
--       { model | status = status }



-- -- VIEW


-- view : Model -> Html Msg
-- view model =
--   div [ class "main" ] [
--     div [ class "signup" ]
--     [ Html.form [ action "http://localhost:4000/v1/todoitems", id "userform", method "POST" ]
--         [ label [ attribute "aria-hidden" "true", for "chk" ]
--             [ text "To-Do List Form" ]
--         , div []
--         [ viewInput "text" "Task Name" model.task_name Task_Name
--         , viewInput "text" "Description" model.description Description
--         , viewInput "text" "Category" model.category Category
--         , viewInput "text" "Priority" model.priority Priority
--         , viewInput "text" "Status" model.status Status
--         , viewValidation model
--         ]
--         , button []
--             [ text "Submit" ]
--         ]
--     ]
--   ]


-- viewInput : String -> String -> String -> (String -> msg) -> Html msg
-- viewInput t p v toMsg =
--   input [ type_ t, placeholder p, value v, onInput toMsg ] []


-- viewValidation : Model -> Html msg
-- viewValidation model =
--   if model.task_name == "" || model.description == "" || model.category == "" || model.priority == "" || model.status == "" then
--     div [ style "color" "red", style "text-align" "center" ] [ text "Please Fill All Fields!" ]
--   else
--     div [ style "color" "green",  style "text-align" "center" ] [ text "Good!" ]