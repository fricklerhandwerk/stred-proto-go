package core

// this package captures the semantics of protobuf documents, disregarding
// presentational matters entierly.

// the design philosphy is to reduce input parameters to identifier strings,
// field numbers, boolean flags, built-in and valid user-defined message types.
// all items have to be created as ancestors of `Document` to be meaningful and
// usable.

// it must be impossible to update the document with an invalid value, ideally
// to even supply an invald value. consequently consumers are not supposed to
// create values of some types other than through due process, which will
// perform proper initialision. these are therefore only visible through the
// interfaces they implement. this is the case for `Message` and `Enum`, which
// can be referenced as RPC message types or message field types; they must
// constitute valid values for that purpose.

// a primary design decision following this principle is that once-valid values
// cannot be invalidated. for example setting an unset field to some valid
// value means that this field's value can only be changed to another valid
// one, but not unset any more.

// the type layout follows a roughly hierarchic pattern: the root object is
// `Document`, from which we can create tentative subsidiary items. for each we
// can set required attributes, and insert the object into the parent it was
// created by.
// tentative objects are not tracked by the document and in case of `NewMessage`
// and `NewEnum` may not have children. first, doing that would complicate the
// necessary type layout to a degree Go does not afford to practically handle
// - a message field would then have a parent of type either `Message` or
// `NewMessage`, but then we would need to add indirection through yet another
// interface.
// second, it prevents a problem with dangling references, since by design the
// document does not account for tentative elements. if a valid tentative item
// is changed to reference a message, and that message is deleted from the
// document, the item gets invalidated inadvertently, violating the stated
// principle of irreversible correctness. the document only cares about
// semantically correct structures, so tentative ones cannot even be referenced
// other than on creation, and cannot be passed into the document through
// exported methods.
// insertion into the parent means hooking up the pointer into a collection in
// the parent, except for `Message` and `Enum`, which are created by copying
// contents. to continue operating on the resulting object in these special
// cases you have to fetch it back from its parent by comparing labels, which
// must be unique. this is not elegant, but preserves a consistent interface.
