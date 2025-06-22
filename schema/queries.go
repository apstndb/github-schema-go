package schema

// Predefined jq queries for common operations

const (
	// typeQuery formats a GraphQL type with all its fields
	typeQuery = `
def formatType:
  if type == "object" and .kind == "NON_NULL" then
    (.ofType | formatType) + "!"
  elif type == "object" and .kind == "LIST" then
    "[" + (.ofType | formatType) + "]"
  elif type == "object" then
    .name // .kind
  else
    .
  end;

.data.__schema.types[] | 
select(.name == $type) |
{
  type: {
    name,
    kind,
    description,
    fields: (
      if .fields then
        [.fields[] | {
          name,
          description,
          type: (.type | formatType),
          arguments: (
            if (.args | length) > 0 then
              [.args[] | {
                name,
                description,
                type: (.type | formatType)
              }]
            else
              null
            end
          )
        }]
      else
        null
      end
    ),
    inputFields: (
      if .inputFields then
        [.inputFields[] | {
          name,
          description,
          type: (.type | formatType),
          required: (.type.kind == "NON_NULL")
        }]
      else
        null
      end
    ),
    enumValues: (
      if .enumValues then
        [.enumValues[] | {
          name,
          description
        }]
      else
        null
      end
    )
  }
}`

	// searchQuery searches for types matching a pattern
	searchQuery = `
[.data.__schema.types[] | 
  select(.name | test($pattern; "i")) | 
  {
    name,
    kind,
    description: (
      if .description != null and (.description | length) > 100 then
        .description[0:100] + "..."
      else
        .description
      end
    )
  }] | {
    count: length,
    pattern: $pattern,
    results: .
  }`

	// mutationQuery formats a mutation with expanded input details
	mutationQuery = `
def formatType:
  if type == "object" and .kind == "NON_NULL" then
    (.ofType | formatType) + "!"
  elif type == "object" and .kind == "LIST" then
    "[" + (.ofType | formatType) + "]"
  elif type == "object" then
    .name // .kind
  else
    .
  end;

# Find the mutation
(.data.__schema.types[] | select(.name == "Mutation").fields[] | select(.name == $mutation)) as $mut |

# Get input type details if it exists  
if $mut.args[0].type.ofType.name then
  (.data.__schema.types[] | select(.name == $mut.args[0].type.ofType.name)) as $inputType |
  {
    mutation: {
      name: $mut.name,
      description: $mut.description,
      inputs: [{
        name: $mut.args[0].name,
        type: ($mut.args[0].type | formatType),
        description: (
          $mut.args[0].description + "\n\nInput object '" + $inputType.name + "' has the following fields:\n" +
          ([$inputType.inputFields[] | 
            "- " + .name + ": " + (.type | formatType) + 
            if .type.kind == "NON_NULL" then " (required)" else "" end +
            if .description then "\n  " + .description else "" end
          ] | join("\n"))
        ),
        required: ($mut.args[0].type.kind == "NON_NULL")
      }]
    }
  }
else
  {
    mutation: {
      name: $mut.name,
      description: $mut.description,
      inputs: [$mut.args[] | {
        name,
        type: (.type | formatType),
        description,
        required: (.type.kind == "NON_NULL")
      }]
    }
  }
end`

	// fieldSearchQuery searches for fields across all types
	fieldSearchQuery = `
[.data.__schema.types[] |
{
  type: .name,
  kind: .kind,
  fields: [.fields[]? | select(.name | test($pattern; "i")) | {
    name,
    type: (
      if .type.kind == "NON_NULL" then
        .type.ofType.name + "!"
      elif .type.kind == "LIST" then
        "[" + (.type.ofType.name // .type.ofType.kind) + "]"
      else
        .type.name
      end
    ),
    description
  }]
} |
select(.fields | length > 0)]`

	// interfaceImplementersQuery finds types implementing an interface
	interfaceImplementersQuery = `
.data.__schema.types[] |
select(.name == $interface) |
if .possibleTypes then
  {
    interface: .name,
    implementers: [.possibleTypes[] | .name]
  }
else
  {
    interface: .name,
    implementers: []
  }
end`
)

// Additional helper queries that can be exposed

const (
	// ListMutationsQuery lists all available mutations
	ListMutationsQuery = `.data.__schema.types[] | select(.name == "Mutation") | .fields[] | .name`

	// ListTypesQuery lists all type names
	ListTypesQuery = `.data.__schema.types[] | .name`

	// ListObjectTypesQuery lists only object types
	ListObjectTypesQuery = `.data.__schema.types[] | select(.kind == "OBJECT") | .name`

	// ListInputTypesQuery lists only input types
	ListInputTypesQuery = `.data.__schema.types[] | select(.kind == "INPUT_OBJECT") | .name`
)