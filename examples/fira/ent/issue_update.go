// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/adnaan/fir/examples/fira/ent/issue"
	"github.com/adnaan/fir/examples/fira/ent/predicate"
	"github.com/adnaan/fir/examples/fira/ent/project"
	"github.com/google/uuid"
)

// IssueUpdate is the builder for updating Issue entities.
type IssueUpdate struct {
	config
	hooks    []Hook
	mutation *IssueMutation
}

// Where appends a list predicates to the IssueUpdate builder.
func (iu *IssueUpdate) Where(ps ...predicate.Issue) *IssueUpdate {
	iu.mutation.Where(ps...)
	return iu
}

// SetUpdateTime sets the "update_time" field.
func (iu *IssueUpdate) SetUpdateTime(t time.Time) *IssueUpdate {
	iu.mutation.SetUpdateTime(t)
	return iu
}

// SetTitle sets the "title" field.
func (iu *IssueUpdate) SetTitle(s string) *IssueUpdate {
	iu.mutation.SetTitle(s)
	return iu
}

// SetDescription sets the "description" field.
func (iu *IssueUpdate) SetDescription(s string) *IssueUpdate {
	iu.mutation.SetDescription(s)
	return iu
}

// SetOwnerID sets the "owner" edge to the Project entity by ID.
func (iu *IssueUpdate) SetOwnerID(id uuid.UUID) *IssueUpdate {
	iu.mutation.SetOwnerID(id)
	return iu
}

// SetNillableOwnerID sets the "owner" edge to the Project entity by ID if the given value is not nil.
func (iu *IssueUpdate) SetNillableOwnerID(id *uuid.UUID) *IssueUpdate {
	if id != nil {
		iu = iu.SetOwnerID(*id)
	}
	return iu
}

// SetOwner sets the "owner" edge to the Project entity.
func (iu *IssueUpdate) SetOwner(p *Project) *IssueUpdate {
	return iu.SetOwnerID(p.ID)
}

// Mutation returns the IssueMutation object of the builder.
func (iu *IssueUpdate) Mutation() *IssueMutation {
	return iu.mutation
}

// ClearOwner clears the "owner" edge to the Project entity.
func (iu *IssueUpdate) ClearOwner() *IssueUpdate {
	iu.mutation.ClearOwner()
	return iu
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (iu *IssueUpdate) Save(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	iu.defaults()
	if len(iu.hooks) == 0 {
		if err = iu.check(); err != nil {
			return 0, err
		}
		affected, err = iu.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*IssueMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = iu.check(); err != nil {
				return 0, err
			}
			iu.mutation = mutation
			affected, err = iu.sqlSave(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(iu.hooks) - 1; i >= 0; i-- {
			if iu.hooks[i] == nil {
				return 0, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = iu.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, iu.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// SaveX is like Save, but panics if an error occurs.
func (iu *IssueUpdate) SaveX(ctx context.Context) int {
	affected, err := iu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (iu *IssueUpdate) Exec(ctx context.Context) error {
	_, err := iu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (iu *IssueUpdate) ExecX(ctx context.Context) {
	if err := iu.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (iu *IssueUpdate) defaults() {
	if _, ok := iu.mutation.UpdateTime(); !ok {
		v := issue.UpdateDefaultUpdateTime()
		iu.mutation.SetUpdateTime(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (iu *IssueUpdate) check() error {
	if v, ok := iu.mutation.Title(); ok {
		if err := issue.TitleValidator(v); err != nil {
			return &ValidationError{Name: "title", err: fmt.Errorf(`ent: validator failed for field "Issue.title": %w`, err)}
		}
	}
	if v, ok := iu.mutation.Description(); ok {
		if err := issue.DescriptionValidator(v); err != nil {
			return &ValidationError{Name: "description", err: fmt.Errorf(`ent: validator failed for field "Issue.description": %w`, err)}
		}
	}
	return nil
}

func (iu *IssueUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   issue.Table,
			Columns: issue.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeUUID,
				Column: issue.FieldID,
			},
		},
	}
	if ps := iu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := iu.mutation.UpdateTime(); ok {
		_spec.SetField(issue.FieldUpdateTime, field.TypeTime, value)
	}
	if value, ok := iu.mutation.Title(); ok {
		_spec.SetField(issue.FieldTitle, field.TypeString, value)
	}
	if value, ok := iu.mutation.Description(); ok {
		_spec.SetField(issue.FieldDescription, field.TypeString, value)
	}
	if iu.mutation.OwnerCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   issue.OwnerTable,
			Columns: []string{issue.OwnerColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUUID,
					Column: project.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := iu.mutation.OwnerIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   issue.OwnerTable,
			Columns: []string{issue.OwnerColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUUID,
					Column: project.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, iu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{issue.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	return n, nil
}

// IssueUpdateOne is the builder for updating a single Issue entity.
type IssueUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *IssueMutation
}

// SetUpdateTime sets the "update_time" field.
func (iuo *IssueUpdateOne) SetUpdateTime(t time.Time) *IssueUpdateOne {
	iuo.mutation.SetUpdateTime(t)
	return iuo
}

// SetTitle sets the "title" field.
func (iuo *IssueUpdateOne) SetTitle(s string) *IssueUpdateOne {
	iuo.mutation.SetTitle(s)
	return iuo
}

// SetDescription sets the "description" field.
func (iuo *IssueUpdateOne) SetDescription(s string) *IssueUpdateOne {
	iuo.mutation.SetDescription(s)
	return iuo
}

// SetOwnerID sets the "owner" edge to the Project entity by ID.
func (iuo *IssueUpdateOne) SetOwnerID(id uuid.UUID) *IssueUpdateOne {
	iuo.mutation.SetOwnerID(id)
	return iuo
}

// SetNillableOwnerID sets the "owner" edge to the Project entity by ID if the given value is not nil.
func (iuo *IssueUpdateOne) SetNillableOwnerID(id *uuid.UUID) *IssueUpdateOne {
	if id != nil {
		iuo = iuo.SetOwnerID(*id)
	}
	return iuo
}

// SetOwner sets the "owner" edge to the Project entity.
func (iuo *IssueUpdateOne) SetOwner(p *Project) *IssueUpdateOne {
	return iuo.SetOwnerID(p.ID)
}

// Mutation returns the IssueMutation object of the builder.
func (iuo *IssueUpdateOne) Mutation() *IssueMutation {
	return iuo.mutation
}

// ClearOwner clears the "owner" edge to the Project entity.
func (iuo *IssueUpdateOne) ClearOwner() *IssueUpdateOne {
	iuo.mutation.ClearOwner()
	return iuo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (iuo *IssueUpdateOne) Select(field string, fields ...string) *IssueUpdateOne {
	iuo.fields = append([]string{field}, fields...)
	return iuo
}

// Save executes the query and returns the updated Issue entity.
func (iuo *IssueUpdateOne) Save(ctx context.Context) (*Issue, error) {
	var (
		err  error
		node *Issue
	)
	iuo.defaults()
	if len(iuo.hooks) == 0 {
		if err = iuo.check(); err != nil {
			return nil, err
		}
		node, err = iuo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*IssueMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = iuo.check(); err != nil {
				return nil, err
			}
			iuo.mutation = mutation
			node, err = iuo.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(iuo.hooks) - 1; i >= 0; i-- {
			if iuo.hooks[i] == nil {
				return nil, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = iuo.hooks[i](mut)
		}
		v, err := mut.Mutate(ctx, iuo.mutation)
		if err != nil {
			return nil, err
		}
		nv, ok := v.(*Issue)
		if !ok {
			return nil, fmt.Errorf("unexpected node type %T returned from IssueMutation", v)
		}
		node = nv
	}
	return node, err
}

// SaveX is like Save, but panics if an error occurs.
func (iuo *IssueUpdateOne) SaveX(ctx context.Context) *Issue {
	node, err := iuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (iuo *IssueUpdateOne) Exec(ctx context.Context) error {
	_, err := iuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (iuo *IssueUpdateOne) ExecX(ctx context.Context) {
	if err := iuo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (iuo *IssueUpdateOne) defaults() {
	if _, ok := iuo.mutation.UpdateTime(); !ok {
		v := issue.UpdateDefaultUpdateTime()
		iuo.mutation.SetUpdateTime(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (iuo *IssueUpdateOne) check() error {
	if v, ok := iuo.mutation.Title(); ok {
		if err := issue.TitleValidator(v); err != nil {
			return &ValidationError{Name: "title", err: fmt.Errorf(`ent: validator failed for field "Issue.title": %w`, err)}
		}
	}
	if v, ok := iuo.mutation.Description(); ok {
		if err := issue.DescriptionValidator(v); err != nil {
			return &ValidationError{Name: "description", err: fmt.Errorf(`ent: validator failed for field "Issue.description": %w`, err)}
		}
	}
	return nil
}

func (iuo *IssueUpdateOne) sqlSave(ctx context.Context) (_node *Issue, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   issue.Table,
			Columns: issue.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeUUID,
				Column: issue.FieldID,
			},
		},
	}
	id, ok := iuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "Issue.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := iuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, issue.FieldID)
		for _, f := range fields {
			if !issue.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != issue.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := iuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := iuo.mutation.UpdateTime(); ok {
		_spec.SetField(issue.FieldUpdateTime, field.TypeTime, value)
	}
	if value, ok := iuo.mutation.Title(); ok {
		_spec.SetField(issue.FieldTitle, field.TypeString, value)
	}
	if value, ok := iuo.mutation.Description(); ok {
		_spec.SetField(issue.FieldDescription, field.TypeString, value)
	}
	if iuo.mutation.OwnerCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   issue.OwnerTable,
			Columns: []string{issue.OwnerColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUUID,
					Column: project.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := iuo.mutation.OwnerIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   issue.OwnerTable,
			Columns: []string{issue.OwnerColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUUID,
					Column: project.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &Issue{config: iuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, iuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{issue.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	return _node, nil
}