<div class="modal-header">
    <h3 class="modal-title">Post Options</h3>
</div>
<div class="modal-body">
    <div class="container-fluid">
        <form class="form-horizontal">
            <div class="form-group">
                <label for="post-slug" class="col-sm-2 control-label">Custom Slug</label>
                <div class="col-sm-4">
                    <input type="text" class="form-control" id="post-slug" ng-model="shared.post.Slug" value="{{shared.post.Slug}}">
                </div>
            </div>
            <div class="form-group">
                <label for="post-cover" class="col-sm-2 control-label">Cover</label>
                <div class="col-sm-10">
                    <a ng-controller="ImageModalCtrl" ng-click="open('lg', 'post-cover')"><img class="img-settings img-thumbnail img-settings" id="post-cover" src="{{shared.post.Image}}" alt="{{shared.post.Image}}" ng-if="shared.post.Image!=''" /><img class="img-settings img-thumbnail img-settings" id="post-cover" src="/public/images/no-image.png" alt="No image" ng-if="shared.post.Image==''" /></a> <a class="text-danger" id="post-cover-delete" ng-controller="EmptyModalCtrl" ng-click="deleteCover()"><span class="glyphicon glyphicon-remove" aria-hidden="true"></span> Remove</a>
                </div>
            </div>
            <div class="form-group">
                <div class="col-sm-offset-1 col-sm-10">
                  <div class="checkbox">
                    <label>
                        <input bs-switch ng-model="shared.post.IsFeatured" type="checkbox" class="post-checkbox" data-label-text="Feature Post" data-label-width="85" data-off-text="NO" data-on-text="YES" data-on-color="success" data-off-color="danger" data-size="normal">
                    </label>
                  </div>
                </div>
            </div>
            <div class="form-group">
                <div class="col-sm-offset-1 col-sm-10">
                  <div class="checkbox">
                    <label>
                        <input bs-switch ng-model="shared.post.IsPage" type="checkbox" class="post-checkbox" data-label-text="Static Page" data-label-width="85" data-off-text="NO" data-on-text="YES" data-on-color="success" data-off-color="danger" data-size="normal">
                    </label>
                  </div>
                </div>
            </div>
        </form>
    </div>
</div>
<div class="modal-footer">
    <button class="btn btn-primary" ng-click="ok()">OK</button>
</div>